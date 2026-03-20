package dto

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/int-xiaoli/go_gateway/public"
	"gopkg.in/go-playground/validator.v9"
	zh_translations "gopkg.in/go-playground/validator.v9/translations/zh"
)

// TestAdminLoginInput_BindValidParam 测试参数绑定和验证的完整流程
func TestAdminLoginInput_BindValidParam(t *testing.T) {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		contentType    string
		body           string
		expectedUser   string
		expectedPass   string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:         "表单数据 - 验证通过",
			contentType:  "application/x-www-form-urlencoded",
			body:         "username=admin&password=12345",
			expectedUser: "admin",
			expectedPass: "12345",
			expectError:  false,
		},
		{
			name:         "JSON数据 - 验证通过",
			contentType:  "application/json",
			body:         `{"username":"admin","password":"12345"}`,
			expectedUser: "admin",
			expectedPass: "12345",
			expectError:  false,
		},
		{
			name:           "表单数据 - 用户名错误",
			contentType:    "application/x-www-form-urlencoded",
			body:           "username=wrong&password=12345",
			expectedUser:   "wrong",
			expectedPass:   "12345",
			expectError:    true,
			expectedErrMsg: "姓名 填写不正确哦",
		},
		{
			name:           "JSON数据 - 缺少密码",
			contentType:    "application/json",
			body:           `{"username":"admin"}`,
			expectedUser:   "admin",
			expectedPass:   "",
			expectError:    true,
			expectedErrMsg: "密码为必填字段",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试请求
			req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", tt.contentType)

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 创建 Gin 上下文
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// 手动设置 validator 和 translator（模拟 TranslationMiddleware 的工作）
			setupValidatorAndTranslator(c)

			// 创建 DTO 实例
			params := &AdminLoginInput{}
			t.Logf("=== 测试开始: %s ===", tt.name)
			t.Logf("输入数据: %s", tt.body)
			t.Logf("DTO 初始状态: UserName='%s', Password='%s'", params.UserName, params.Password)

			// 执行绑定和验证
			err := params.BindValidParam(c)

			// 检查结果
			t.Logf("DTO 绑定后状态: UserName='%s', Password='%s'", params.UserName, params.Password)

			if err != nil {
				t.Logf("验证错误: %v", err.Error())
			} else {
				t.Logf("验证成功")
			}

			// 验证绑定结果
			if params.UserName != tt.expectedUser {
				t.Errorf("期望 UserName='%s', 实际='%s'", tt.expectedUser, params.UserName)
			}
			if params.Password != tt.expectedPass {
				t.Errorf("期望 Password='%s', 实际='%s'", tt.expectedPass, params.Password)
			}

			// 验证错误处理
			if tt.expectError {
				if err == nil {
					t.Errorf("期望错误，但没有发生错误")
				} else if tt.expectedErrMsg != "" && tt.expectedErrMsg != err.Error() {
					t.Errorf("期望错误信息='%s', 实际='%s'", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("不期望错误，但发生了错误: %v", err)
				}
			}

			t.Logf("=== 测试结束: %s ===\n", tt.name)
		})
	}
}

// setupValidatorAndTranslator 模拟 TranslationMiddleware 中的 validator 和 translator 设置
func setupValidatorAndTranslator(c *gin.Context) {
	// 设置中文翻译器
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")

	// 创建 validator 实例
	val := validator.New()

	// 注册默认中文翻译
	zh_translations.RegisterDefaultTranslations(val, trans)

	// 注册字段名映射函数
	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get("comment")
	})

	// 注册自定义验证规则
	val.RegisterValidation("is_valid_username", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		// 添加调试日志
		println("执行自定义验证: 字段值 =", value, ", 期望值 = admin")
		return value == "admin"
	})

	// 注册自定义翻译器
	val.RegisterTranslation("is_valid_username", trans,
		func(ut ut.Translator) error {
			return ut.Add("is_valid_username", "{0} 填写不正确哦", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("is_valid_username", fe.Field())
			return t
		})

	// 存储到 Context
	c.Set(public.TranslatorKey, trans)
	c.Set(public.ValidatorKey, val)
}

// TestManualBinding 手动测试绑定过程，更详细地展示每个步骤
func TestManualBinding(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建测试数据
	testData := map[string]interface{}{
		"username": "admin",
		"password": "12345",
	}
	jsonData, _ := json.Marshal(testData)

	// 创建请求和上下文
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// 设置验证器
	setupValidatorAndTranslator(c)

	// 步骤1: 创建空的 DTO
	params := &AdminLoginInput{}
	t.Logf("步骤1 - DTO 初始状态: %+v", params)

	// 步骤2: 手动执行 ShouldBind
	err := c.ShouldBind(params)
	if err != nil {
		t.Fatalf("ShouldBind 失败: %v", err)
	}
	t.Logf("步骤2 - ShouldBind 后: %+v", params)

	// 步骤3: 手动执行验证
	valid, _ := public.GetValidator(c)
	err = valid.Struct(params)
	if err != nil {
		t.Logf("步骤3 - 验证失败: %v", err)
	} else {
		t.Logf("步骤3 - 验证成功")
	}

	// 步骤4: 如果有错误，进行翻译
	if err != nil {
		trans, _ := public.GetTranslation(c)
		validationErrs := err.(validator.ValidationErrors)
		for _, e := range validationErrs {
			translated := e.Translate(trans)
			t.Logf("步骤4 - 翻译后的错误: %s", translated)
		}
	}
}
