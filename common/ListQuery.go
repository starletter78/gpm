package common

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"gpm/global"

	"gorm.io/gorm"
)

// 错误变量定义（中文提示，支持国际化扩展）
var (
	ErrTargetNotSlicePtr    = errors.New("目标必须是切片指针")
	ErrSourceNotStruct      = errors.New("源类型必须是结构体或结构体指针")
	ErrTargetNotStruct      = errors.New("目标类型必须是结构体或结构体指针")
	ErrSortFieldNotAllowed  = errors.New("排序字段不允许")
	ErrJoinConditionInvalid = errors.New("联表条件无效")
)

// PageInfo 分页查询基本参数
type PageInfo struct {
	Limit int    `form:"limit"` // 每页条数
	Page  int    `form:"page"`  // 页码
	Key   string `form:"key"`   // 搜索关键词
	Order string `form:"order"` // 排序字段，格式: 字段1:asc,字段2:desc
}

// GetPage 获取安全的页码（限制范围）
func (p PageInfo) GetPage() int {
	if p.Page > 20 || p.Page <= 0 {
		return 1
	}
	return p.Page
}

// GetLimit 获取安全的每页条数（限制范围）
func (p PageInfo) GetLimit() int {
	if p.Limit <= 0 || p.Limit > 50 {
		return 10
	}
	return p.Limit
}

// GetOffset 计算偏移量
func (p PageInfo) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

// JoinCondition 联表查询条件（支持参数化，增强安全性）
type JoinCondition struct {
	JoinType string        // 连接类型：INNER, LEFT, RIGHT, FULL等
	Table    string        // 关联表名（支持别名，如"users u"）
	On       string        // 连接条件，支持占位符，如"u.id = t.user_id AND u.status = ?"
	Args     []interface{} // 条件参数，与On中的占位符对应
}

// SortField 排序字段配置
type SortField struct {
	Field     string // 排序字段名
	Direction string // 排序方向: asc/desc
}

// Options 列表查询选项（扩展默认排序配置）
type Options struct {
	PageInfo         PageInfo        // 分页信息
	Likes            []string        // 模糊查询字段
	Preloads         []string        // 预加载关联（ORM层面）
	Joins            []JoinCondition // 联表查询条件（SQL层面）
	Where            *gorm.DB        // 自定义查询条件
	DefaultOrder     string          // 默认排序，格式: 字段1:asc,字段2:desc（为空则不设置默认排序）
	AllowedSorts     []string        // 允许排序的字段列表，用于防注入
	SelectFields     []string        // 需要保留的字段列表（为空则查询所有字段）
	OmitFields       []string        // 需要排除的字段列表
	Context          context.Context // 上下文
	EnableFieldCache bool            // 是否启用字段映射缓存（提升性能）
}

// QueryBuilder 通用查询构建器
type QueryBuilder[T any] struct {
	model   T        // 数据模型
	options Options  // 查询选项
	query   *gorm.DB // GORM查询对象
	count   int64    // 总记录数
	err     error    // 错误信息
}

// 字段映射缓存（键：源类型+目标类型，值：字段映射关系）
var fieldMapCache = sync.Map{}

// NewQueryBuilder 创建新的查询构建器
func NewQueryBuilder[T any](model T, options Options) *QueryBuilder[T] {
	// 初始化默认值
	if options.Context == nil {
		options.Context = context.Background()
	}
	// 默认启用缓存
	if options.EnableFieldCache {
		options.EnableFieldCache = true
	}
	return &QueryBuilder[T]{
		model:   model,
		options: options,
	}
}

// Build 构建完整查询
func (qb *QueryBuilder[T]) Build() *QueryBuilder[T] {
	// 初始化基础查询
	qb.initBaseQuery()
	if qb.err != nil {
		return qb
	}

	// 应用字段筛选
	qb.applyFieldFilters()

	// 处理联表查询（参数化，增强安全性）
	qb.applyJoins()

	// 应用模糊查询
	qb.applyLikeSearch()

	// 应用自定义查询条件
	qb.applyCustomWhere()

	// 预加载关联
	qb.applyPreloads()

	// 计算总记录数
	qb.calculateTotalCount()
	if qb.err != nil {
		return qb
	}

	// 应用分页
	qb.applyPagination()

	// 应用排序
	qb.applySorting()

	return qb
}

// 初始化基础查询
func (qb *QueryBuilder[T]) initBaseQuery() {
	qb.query = global.DB.WithContext(qb.options.Context).Model(qb.model).Where(qb.model)
}

// 应用字段筛选（选择/排除字段）
func (qb *QueryBuilder[T]) applyFieldFilters() {
	if qb.err != nil {
		return
	}

	// 优先处理保留字段
	if len(qb.options.SelectFields) > 0 {
		qb.query = qb.query.Select(qb.options.SelectFields)
		return
	}

	// 处理排除字段
	if len(qb.options.OmitFields) > 0 {
		qb.query = qb.query.Omit(qb.options.OmitFields...)
	}
}

// 应用联表查询条件（参数化处理，防止注入）
func (qb *QueryBuilder[T]) applyJoins() {
	if qb.err != nil || len(qb.options.Joins) == 0 {
		return
	}

	for _, join := range qb.options.Joins {
		// 验证表名格式（简单防护）
		if !isValidTableName(join.Table) {
			qb.err = ErrJoinConditionInvalid
			return
		}
		// 构建参数化联表语句
		joinStr := fmt.Sprintf("%s JOIN %s ON %s", join.JoinType, join.Table, join.On)
		qb.query = qb.query.Joins(joinStr, join.Args...)
	}
}

// 验证表名合法性（防止注入）
func isValidTableName(table string) bool {
	// 允许字母、数字、下划线和空格（别名场景）
	for _, c := range table {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') || c == '_' || c == ' ') {
			return false
		}
	}
	return true
}

// 应用模糊搜索
func (qb *QueryBuilder[T]) applyLikeSearch() {
	if qb.err != nil || len(qb.options.Likes) == 0 || qb.options.PageInfo.Key == "" {
		return
	}

	likes := global.DB.Where("")
	for _, column := range qb.options.Likes {
		likes.Or(
			fmt.Sprintf("%s like ?", column),
			fmt.Sprintf("%%%s%%", qb.options.PageInfo.Key))
	}
	qb.query = qb.query.Where(likes)
}

// 应用自定义查询条件
func (qb *QueryBuilder[T]) applyCustomWhere() {
	if qb.err != nil || qb.options.Where == nil {
		return
	}
	qb.query = qb.query.Where(qb.options.Where)
}

// 应用预加载关联
func (qb *QueryBuilder[T]) applyPreloads() {
	if qb.err != nil || len(qb.options.Preloads) == 0 {
		return
	}

	for _, preload := range qb.options.Preloads {
		qb.query = qb.query.Preload(preload)
	}
}

// 计算总记录数
func (qb *QueryBuilder[T]) calculateTotalCount() {
	if qb.err != nil {
		return
	}

	if err := qb.query.Count(&qb.count).Error; err != nil {
		qb.err = fmt.Errorf("计算总记录数失败: %w", err)
	}
}

// 应用分页
func (qb *QueryBuilder[T]) applyPagination() {
	if qb.err != nil {
		return
	}

	limit := qb.options.PageInfo.GetLimit()
	offset := qb.options.PageInfo.GetOffset()
	qb.query = qb.query.Offset(offset).Limit(limit)
}

// 应用排序
func (qb *QueryBuilder[T]) applySorting() {
	if qb.err != nil {
		return
	}

	var sorts []SortField
	var err error

	// 优先使用用户指定的排序
	if qb.options.PageInfo.Order != "" {
		sorts, err = qb.parseSort(qb.options.PageInfo.Order, qb.options.AllowedSorts)
	} else if qb.options.DefaultOrder != "" {
		// 使用默认排序（支持空，为空则不设置默认排序）
		sorts, err = qb.parseSort(qb.options.DefaultOrder, []string{})
	}

	if err != nil {
		qb.err = err
		return
	}

	// 应用排序
	if len(sorts) > 0 {
		qb.query = qb.applySorts(sorts)
	}
}

// 解析排序字符串为SortField列表
func (qb *QueryBuilder[T]) parseSort(sortStr string, allowedSorts []string) ([]SortField, error) {
	var sorts []SortField
	if sortStr == "" {
		return sorts, nil
	}

	parts := strings.Split(sortStr, ",")
	for _, part := range parts {
		fieldDir := strings.Split(strings.TrimSpace(part), ":")
		if len(fieldDir) < 1 {
			continue
		}

		field := strings.TrimSpace(fieldDir[0])
		// 检查排序字段是否在允许的列表中
		if len(allowedSorts) > 0 {
			allowed := false
			for _, allow := range allowedSorts {
				if field == allow {
					allowed = true
					break
				}
			}
			if !allowed {
				return nil, ErrSortFieldNotAllowed
			}
		}

		// 处理排序方向
		direction := "asc"
		if len(fieldDir) > 1 {
			dir := strings.ToLower(strings.TrimSpace(fieldDir[1]))
			if dir == "desc" {
				direction = "desc"
			}
		}

		sorts = append(sorts, SortField{
			Field:     field,
			Direction: direction,
		})
	}

	return sorts, nil
}

// 将排序条件应用到查询
func (qb *QueryBuilder[T]) applySorts(sorts []SortField) *gorm.DB {
	query := qb.query
	for _, sort := range sorts {
		query = query.Order(fmt.Sprintf("%s %s", sort.Field, sort.Direction))
	}
	return query
}

// GetResult 执行查询并获取结果
func (qb *QueryBuilder[T]) GetResult() (list []T, count int64, err error) {
	if qb.err != nil {
		return nil, 0, qb.err
	}
	if qb.query == nil {
		return nil, 0, errors.New("查询构建器未初始化")
	}
	if err := qb.query.Find(&list).Error; err != nil {
		return nil, 0, fmt.Errorf("执行查询失败: %w", err)
	}

	return list, qb.count, nil
}

// MapToTarget 将查询结果映射到目标结构体切片（支持类型转换和嵌套结构体）
func (qb *QueryBuilder[T]) MapToTarget(target interface{}, excludeFields []string) (interface{}, int64, error) {
	sourceList, count, err := qb.GetResult()
	if err != nil {
		return nil, 0, err
	}

	// 验证目标参数必须是切片指针
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.Elem().Kind() != reflect.Slice {
		return nil, 0, ErrTargetNotSlicePtr
	}

	// 执行类型映射，显式指定泛型类型参数
	mappedResult, err := mapTo[T](sourceList, target, excludeFields, qb.options.EnableFieldCache)
	if err != nil {
		return nil, 0, err
	}

	return mappedResult, count, nil
}

// mapTo 通用类型映射函数（支持类型转换、嵌套结构体、缓存）
func mapTo[Source any](sourceList []Source, target interface{}, excludeFields []string, enableCache bool) (interface{}, error) {
	targetVal := reflect.ValueOf(target)
	if targetVal.Kind() != reflect.Ptr || targetVal.Elem().Kind() != reflect.Slice {
		return nil, ErrTargetNotSlicePtr
	}

	// 获取目标切片的元素类型
	elemType := targetVal.Elem().Type().Elem()

	// 创建结果切片
	resultSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(sourceList))
	if len(sourceList) == 0 {
		return resultSlice.Interface(), nil
	}

	// 处理排除字段
	excludeSet := make(map[string]bool)
	for _, field := range excludeFields {
		excludeSet[strings.ToLower(field)] = true
	}

	// 获取源类型信息
	sourceType := reflect.TypeOf(*new(Source))
	if sourceType.Kind() == reflect.Ptr {
		sourceType = sourceType.Elem()
	}
	if sourceType.Kind() != reflect.Struct {
		return nil, ErrSourceNotStruct
	}

	// 处理目标类型
	targetElemType := elemType
	if targetElemType.Kind() == reflect.Ptr {
		targetElemType = targetElemType.Elem()
	}
	if targetElemType.Kind() != reflect.Struct {
		return nil, ErrTargetNotStruct
	}

	// 获取字段映射关系（从缓存或新建）
	var fieldMap map[string]int
	cacheKey := fmt.Sprintf("%s->%s", sourceType.String(), targetElemType.String())

	if enableCache {
		if cached, ok := fieldMapCache.Load(cacheKey); ok {
			fieldMap = cached.(map[string]int)
		} else {
			fieldMap = buildFieldMap(targetElemType, excludeSet)
			fieldMapCache.Store(cacheKey, fieldMap)
		}
	} else {
		fieldMap = buildFieldMap(targetElemType, excludeSet)
	}

	// 映射每个元素
	for _, source := range sourceList {
		sourceVal := reflect.ValueOf(source)
		if sourceVal.Kind() == reflect.Ptr {
			sourceVal = sourceVal.Elem()
		}

		// 创建目标元素
		targetElem := reflect.New(targetElemType).Elem()

		// 复制匹配的字段值（支持嵌套结构体和类型转换）
		copyFields(sourceVal, targetElem, fieldMap, excludeSet, enableCache)

		resultSlice = reflect.Append(resultSlice, targetElem)
	}

	return resultSlice.Interface(), nil
}

// 构建字段映射关系（目标字段名到索引）
func buildFieldMap(targetType reflect.Type, excludeSet map[string]bool) map[string]int {
	fieldMap := make(map[string]int)
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		fieldName := strings.ToLower(field.Name)
		if !excludeSet[fieldName] {
			fieldMap[fieldName] = i
		}
	}
	return fieldMap
}

// 复制字段值（支持嵌套结构体和类型转换）
func copyFields(sourceVal, targetVal reflect.Value, fieldMap map[string]int, excludeSet map[string]bool, enableCache bool) {
	sourceType := sourceVal.Type()

	for i := 0; i < sourceType.NumField(); i++ {
		sourceField := sourceType.Field(i)
		fieldName := strings.ToLower(sourceField.Name)

		// 跳过需要排除的字段
		if excludeSet[fieldName] {
			continue
		}

		// 检查目标结构体是否有对应字段
		targetFieldIdx, exists := fieldMap[fieldName]
		if !exists {
			continue
		}

		// 获取源字段和目标字段的值
		sourceFieldVal := sourceVal.Field(i)
		targetField := targetVal.Field(targetFieldIdx)

		// 处理嵌套结构体
		if sourceFieldVal.Kind() == reflect.Struct && targetField.Kind() == reflect.Struct &&
			// 使用反射类型比较，修复类型比较错误
			sourceFieldVal.Type() != reflect.TypeOf(time.Time{}) { // 排除time.Time
			// 递归复制嵌套结构体字段
			nestedFieldMap := buildFieldMap(targetField.Type(), excludeSet)
			copyFields(sourceFieldVal, targetField, nestedFieldMap, excludeSet, enableCache)
			continue
		}

		// 尝试赋值（支持类型转换）
		if !setFieldValue(sourceFieldVal, targetField) {
			// 类型不匹配且无法转换，跳过该字段
			continue
		}
	}
}

// 设置字段值（支持常见类型转换）
func setFieldValue(source, target reflect.Value) bool {
	// 直接赋值（类型匹配）
	if source.Type().AssignableTo(target.Type()) {
		target.Set(source)
		return true
	}

	// 类型转换：int系列 -> string
	if (source.Kind() >= reflect.Int && source.Kind() <= reflect.Uint64) && target.Kind() == reflect.String {
		target.SetString(fmt.Sprintf("%v", source.Interface()))
		return true
	}

	// 类型转换：float系列 -> string
	if (source.Kind() == reflect.Float32 || source.Kind() == reflect.Float64) && target.Kind() == reflect.String {
		target.SetString(fmt.Sprintf("%v", source.Interface()))
		return true
	}

	// 类型转换：time.Time -> string (RFC3339格式)
	if source.Type() == reflect.TypeOf(time.Time{}) && target.Kind() == reflect.String {
		target.SetString(source.Interface().(time.Time).Format(time.RFC3339))
		return true
	}

	// 类型转换：string -> int系列
	if source.Kind() == reflect.String && (target.Kind() >= reflect.Int && target.Kind() <= reflect.Uint64) {
		// 实际项目中可根据需要扩展字符串转数字的逻辑
		return false
	}

	return false
}
