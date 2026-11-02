@ECHO OFF
SETLOCAL
SET SCRIPT_DIR=%~dp0
SET MAVEN_PROJECTBASEDIR=%SCRIPT_DIR%
SET WRAPPER_JAR=%SCRIPT_DIR%.mvn\wrapper\maven-wrapper.jar
SET WRAPPER_PROPS=%SCRIPT_DIR%.mvn\wrapper\maven-wrapper.properties
IF NOT EXIST "%WRAPPER_JAR%" (
  powershell -NoProfile -Command "Get-Content '%WRAPPER_PROPS%' | ForEach-Object { if ($_ -match 'wrapperUrl=(.*)') { $url=$Matches[1]; Invoke-WebRequest -Uri $url -OutFile '%WRAPPER_JAR%' } }"
)
java -Djavax.net.ssl.trustStoreType=Windows-ROOT -classpath "%WRAPPER_JAR%" "-Dmaven.multiModuleProjectDirectory=%MAVEN_PROJECTBASEDIR%" org.apache.maven.wrapper.MavenWrapperMain %*
EXIT /B %ERRORLEVEL%
