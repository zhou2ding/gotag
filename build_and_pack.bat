@echo off
go mod tidy

cd /d %~dp0
set cwd=%cd%

set BuildDate=%date:~0,4%-%date:~5,2%-%date:~8,2%
set PackName=gotag-%BuildDate%

if not exist .\output (
    mkdir .\output
) else (
	if exist .\output\%PackName% (
		rd /S /Q .\output\%PackName%
	)
)

mkdir .\output\%PackName%
mkdir .\output\%PackName%\bin
mkdir .\output\%PackName%\config
mkdir .\output\%PackName%\data
mkdir .\output\%PackName%\data\sample

copy .\config\* %cwd%\output\%PackName%\config
copy .\data\sample\* %cwd%\output\%PackName%\data\sample

cd .\cmd
cd .\sample\
go build -o sample.exe main.go
move .\sample.exe %cwd%\output\%PackName%\bin

cd ..\..\

pause
