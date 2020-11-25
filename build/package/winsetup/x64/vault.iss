#include "..\env.iss"

#define AppName    "vault"
#define AppVersion "0.0.0"

[Setup]
AppName={#AppName}
AppVersion={#AppVersion}
AppVerName={#AppName} {#AppVersion}

ArchitecturesAllowed=x64
ArchitecturesInstallIn64BitMode=x64
ChangesEnvironment=yes

DefaultDirName={autopf}\{#AppName}

OutputBaseFilename=setup_{#AppName}_{#AppVersion}_x64
OutputDir="..\..\..\..\dist"

[Files]
Source: "..\..\..\..\vault.exe"; DestDir: "{app}\bin"

[Code]
procedure CurStepChanged(CurStep: TSetupStep);
begin
    if CurStep = ssPostInstall 
    then EnvAddPath(ExpandConstant('{app}\bin'));
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
    if CurUninstallStep = usPostUninstall 
    then EnvRemovePath(ExpandConstant('{app}\bin'));
end;
