@echo on
cd /d "%~dp0"
echo Iniciando servidor JSON simulado...
start /min json-server --watch usuarios.json --port 3000
exit
