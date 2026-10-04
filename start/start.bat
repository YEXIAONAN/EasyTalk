@echo off
rem EasyTalk 一键启动脚本（Windows，双击运行）
setlocal
chcp 65001 >nul

rem 脚本所在目录与项目根目录
set "SCRIPT_DIR=%~dp0"
set "ROOT_DIR=%SCRIPT_DIR%.."

cd /d "%ROOT_DIR%"

rem 首次运行：从示例创建 config.json
if not exist config.json (
    if exist config.example.json (
        copy /y config.example.json config.json >nul
        echo [EasyTalk] 已从 config.example.json 创建 config.json，请填入你的 API Key。
    ) else (
        echo [EasyTalk] 未找到 config.example.json，请手动创建 config.json。
        exit /b 1
    )
)

rem 首次运行：从示例创建 .env（可选）
if not exist .env (
    if exist .env.example (
        copy /y .env.example .env >nul
        echo [EasyTalk] 已从 .env.example 创建 .env。
    )
)

rem 定位可执行文件
if exist easytalk.exe (
    easytalk.exe
) else if exist easytalk (
    easytalk
) else (
    echo [EasyTalk] 未找到 easytalk 可执行文件，请先构建或下载 Release。
    exit /b 1
)