@echo off
set PATH=%~dp0..\..\scripts;%PATH%
cd %~dp0

call python main_dll.py

pause
