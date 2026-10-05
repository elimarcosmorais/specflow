Unicode true

####
## Instalador do Dev Tools (spec v1, secao 3.3.4; ADR 003 - aplicacao portatil).
## Gerado pelo Wails e ajustado: roda como usuario comum, sugere C:\dev-tools, nao escreve no
## registro e tem um desinstalador que nunca toca em $INSTDIR\projects.
##
## Para depurar o instalador sem o Wails, gere antes o binario com `wails build` e rode:
## > makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\dev-tools.exe project.nsi
####
!define INFO_PROJECTNAME    "dev-tools"
!define INFO_COMPANYNAME    "Dev Tools"
!define INFO_PRODUCTNAME    "Dev Tools"
!define INFO_PRODUCTVERSION "0.1.0"
!define INFO_COPYRIGHT      ""

# Nunca "admin" (ADR 003).
!define REQUEST_EXECUTION_LEVEL "user"

!include "wails_tools.nsh"

# O $INSTDIR do desinstalador e a pasta do uninstall.exe; precisamos do pai dela.
!insertmacro un.GetParent

# The version information for this two must consist of 4 parts
VIProductVersion "${INFO_PRODUCTVERSION}.0"
VIFileVersion    "${INFO_PRODUCTVERSION}.0"

VIAddVersionKey "CompanyName"     "${INFO_COMPANYNAME}"
VIAddVersionKey "FileDescription" "${INFO_PRODUCTNAME} Installer"
VIAddVersionKey "ProductVersion"  "${INFO_PRODUCTVERSION}"
VIAddVersionKey "FileVersion"     "${INFO_PRODUCTVERSION}"
VIAddVersionKey "LegalCopyright"  "${INFO_COPYRIGHT}"
VIAddVersionKey "ProductName"     "${INFO_PRODUCTNAME}"

# Enable HiDPI support. https://nsis.sourceforge.io/Reference/ManifestDPIAware
ManifestDPIAware true

!include "MUI.nsh"

!define MUI_ICON "..\icon.ico"
!define MUI_UNICON "..\icon.ico"
!define MUI_FINISHPAGE_NOAUTOCLOSE # Wait on the INSTFILES page so the user can take a look into the details of the installation steps
!define MUI_ABORTWARNING # This will warn the user if they exit from the installer.

!insertmacro MUI_PAGE_WELCOME # Welcome to the installer page.
!insertmacro MUI_PAGE_DIRECTORY # In which folder install page.
!insertmacro MUI_PAGE_INSTFILES # Installing page.
!insertmacro MUI_PAGE_FINISH # Finished installation page.

!insertmacro MUI_UNPAGE_INSTFILES # Uinstalling page

!insertmacro MUI_LANGUAGE "English"
!insertmacro MUI_LANGUAGE "PortugueseBR"

LangString MSG_REMOVE_DATA ${LANG_ENGLISH}     "Also remove the Dev Tools data and logs?"
LangString MSG_REMOVE_DATA ${LANG_PORTUGUESEBR} "Remover também os dados e os logs do Dev Tools?"

Name "${INFO_PRODUCTNAME}"
OutFile "..\..\bin\${INFO_PROJECTNAME}-${ARCH}-installer.exe" # Name of the installer's file.
InstallDir "C:\dev-tools" # Pasta sugerida; o usuario pode trocar na pagina de pasta.
ShowInstDetails show # This will always show the installation details.

Function .onInit
   !insertmacro wails.checkArchitecture
FunctionEnd

# $INSTDIR passa a ser a raiz da instalacao (a pasta acima de programs\), como no instalador.
Function un.onInit
    ${un.GetParent} "$INSTDIR" $INSTDIR
FunctionEnd

Section
    SetShellVarContext current

    !insertmacro wails.webview2runtime

    CreateDirectory "$INSTDIR\programs\data"
    CreateDirectory "$INSTDIR\programs\logs"
    CreateDirectory "$INSTDIR\projects\web"
    CreateDirectory "$INSTDIR\projects\android"
    CreateDirectory "$INSTDIR\projects\windows"

    SetOutPath "$INSTDIR\programs"

    !insertmacro wails.files

    CreateShortcut "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\programs\${PRODUCT_EXECUTABLE}"
    CreateShortCut "$DESKTOP\${INFO_PRODUCTNAME}.lnk" "$INSTDIR\programs\${PRODUCT_EXECUTABLE}"
    CreateShortcut "$SMPROGRAMS\Desinstalar ${INFO_PRODUCTNAME}.lnk" "$INSTDIR\programs\uninstall.exe"

    WriteUninstaller "$INSTDIR\programs\uninstall.exe"
SectionEnd

Section "uninstall"
    SetShellVarContext current

    # Nenhum comando abaixo toca em $INSTDIR\projects.
    # /SD IDNO: no modo silencioso (/S) os dados e os logs ficam.
    MessageBox MB_YESNO|MB_ICONQUESTION "$(MSG_REMOVE_DATA)" /SD IDNO IDYES removeAll

    # Nao: remove so o programa; data\ e logs\ ficam.
    Delete "$INSTDIR\programs\${PRODUCT_EXECUTABLE}"
    Delete "$INSTDIR\programs\uninstall.exe"
    RMDir /r "$INSTDIR\programs\webview"
    Goto removeShortcuts

    removeAll:
    # Sim: remove programs\ inteira (programa, data\, logs\ e webview\).
    RMDir /r "$INSTDIR\programs"

    removeShortcuts:
    Delete "$SMPROGRAMS\${INFO_PRODUCTNAME}.lnk"
    Delete "$DESKTOP\${INFO_PRODUCTNAME}.lnk"
    Delete "$SMPROGRAMS\Desinstalar ${INFO_PRODUCTNAME}.lnk"

    # Sem /r: so removem pastas vazias.
    RMDir "$INSTDIR\programs"
    RMDir "$INSTDIR"
SectionEnd
