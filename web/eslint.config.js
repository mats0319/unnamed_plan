import stylistic from "@stylistic/eslint-plugin" // eslint将格式化的内容提取到这个包维护了
import tsEslint from "typescript-eslint" // 比原生的defineConfig更智能
import vueEslintParser from "vue-eslint-parser" // 解析vue代码
import pluginVue from "eslint-plugin-vue" // 用于特殊条件（例如本文件中的vue html部分缩进控制）

export default tsEslint.config(
    { ignores: [ "node_modules/**", "dist/**", "public/**", "format_result.html" ] }, // 全局忽略
    // 没有引入任何其他配置，保证eslint不会执行任何非预期行为
    {
        files: [ "**/*.{js,ts,vue}" ],
        plugins: {
            "@stylistic": stylistic,
            "vue": pluginVue,
        },
        languageOptions: {
            parser: vueEslintParser, // 解析vue html
            parserOptions: {
                parser: tsEslint.parser, // 解析vue script(ts)
                ecmaVersion: "latest",
                sourceType: "module",
                extraFileExtensions: [ ".vue" ],
            },
        },
        rules: {
            // style, stylistic不推荐直接启用所有规则并应用其默认值，所以我们一个一个过
            "@stylistic/array-bracket-spacing": [ "warn", "always", { objectsInArrays: false, arraysInArrays: false }], // 数组括号间距
            "@stylistic/arrow-spacing": "warn", // 箭头符号左右应有空格
            "@stylistic/block-spacing": "warn", // 块间距
            "@stylistic/comma-dangle": [ "warn", "only-multiline" ], // 对象和数组字面量的尾随逗号
            "@stylistic/comma-spacing": "warn", // 逗号后应有空格
            "@stylistic/dot-location": "warn", // 链式调用，点和前面的部分在一行，例如`res.err`，如果需要换行应写成`res.\nerr`
            "@stylistic/eol-last": "warn", // 文件末尾应有换行符
            "@stylistic/function-call-spacing": "warn", // 函数调用，函数名和括号中间不应有空格
            "@stylistic/indent": [ "warn", 4 ], // 缩进，4个空格
            "@stylistic/indent-binary-ops": [ "warn", 4 ], // 多行二元运算符缩进，4个空格（推荐与上一条一起使用?）
            "vue/html-indent": [ "warn", 2 ], // vue html部分缩进，2个空格
            "@stylistic/key-spacing": "warn", // 冒号后应有空格
            "@stylistic/keyword-spacing": "warn", // 关键字前后应有空格
            "@stylistic/max-len": [ "warn", {
                code: 120,
                ignoreComments: true,
                ignoreTrailingComments: true,
                ignoreUrls: true,
            }], // 单行代码长度
            "@stylistic/no-multi-spaces": "warn", // 禁止连续空格
            "@stylistic/no-multiple-empty-lines": "warn", // 禁止多空行
            "@stylistic/no-trailing-spaces": "warn", // 禁止行末空格
            "@stylistic/object-curly-spacing": [ "warn", "always" ], // 大括号内部应有空格
            "@stylistic/semi": [ "warn", "never" ], // 分号
            "@stylistic/space-before-blocks": "warn", // 块前空格
            "@stylistic/space-before-function-paren": [ "warn", "never" ], // 函数定义，函数名和括号中间不应有空格
            "@stylistic/space-in-parens": "warn", // 括号里侧不应有空格
            "@stylistic/space-infix-ops": [ "warn", { ignoreTypes: true }], // 中缀运算符前后应有空格，例如`+`/`=`
            "@stylistic/spaced-comment": "warn", // 注释符号和正文中间应有空格

            // eslint
            "no-duplicate-imports": [ "warn", { includeExports: true }], // 一个文件一行导入
            "no-var": "warn", // 使用let/const代替var
            "prefer-const": [ "warn", { destructuring: "all" }], // 优先定义常量
        },
    },
)
