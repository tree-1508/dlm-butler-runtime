param(
    [Parameter(Mandatory = $true)]
    [string]$BinaryPath,
    [Parameter(Mandatory = $true)]
    [string]$OutDir
)

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$binary = (Resolve-Path -LiteralPath $BinaryPath).Path
New-Item -ItemType Directory -Force -Path $OutDir | Out-Null
$licensesDir = Join-Path $OutDir 'licenses'
New-Item -ItemType Directory -Force -Path $licensesDir | Out-Null

# Security #150 comment 5465627296 Q3(B) permits exact-revision-bound,
# authoritative public developer/module metadata as license evidence.
# The allowlist is intentionally exact module+version keyed: no organization-level inference.
$curatedPath = '.github/compliance/dld051-curated-license-evidence.json'
$curated = @{}
if (Test-Path -LiteralPath $curatedPath) {
    $curatedDoc = Get-Content -LiteralPath $curatedPath -Raw | ConvertFrom-Json
    foreach ($entry in @($curatedDoc.entries)) {
        $key = "$($entry.module)@$($entry.version)"
        if ($curated.ContainsKey($key)) { throw "Duplicate curated license evidence key: $key" }
        $curated[$key] = $entry
    }
    Copy-Item -LiteralPath $curatedPath -Destination (Join-Path $OutDir 'CURATED_LICENSE_EVIDENCE.json') -Force
}

# The executable is the source of truth for the exact linked Go module set.
$buildInfo = (& go version -m $binary 2>&1 | Out-String)
if ($LASTEXITCODE -ne 0) { throw 'go version -m failed for the exact butler.exe' }
$buildInfo | Set-Content -Encoding utf8 (Join-Path $OutDir 'butler-go-version-m.txt')

$deps = @()
foreach ($line in ($buildInfo -split "`r?`n")) {
    if ($line -match '^\s*dep\s+([^\s]+)\s+([^\s]+)(?:\s+(h1:[^\s]+))?') {
        $deps += [pscustomobject]@{
            Path = $Matches[1]
            Version = $Matches[2]
            Sum = if ($Matches.Count -ge 4) { $Matches[3] } else { '' }
        }
    }
}
if ($deps.Count -eq 0) { throw 'No linked Go modules were found in butler.exe' }

$moduleTsv = @('module`tversion`tgo_sum')
foreach ($dep in $deps) {
    $moduleTsv += "$($dep.Path)`t$($dep.Version)`t$($dep.Sum)"
}
$moduleTsv | Set-Content -Encoding utf8 (Join-Path $OutDir 'EXACT_BINARY_MODULES.tsv')

$notice = New-Object System.Text.StringBuilder
[void]$notice.AppendLine('DLD Butler Runtime - exact binary third-party notices')
[void]$notice.AppendLine('Generated from: go version -m <exact butler.exe>')
[void]$notice.AppendLine('The module list below contains only modules embedded in the distributed Go executable.')
[void]$notice.AppendLine('Exact-revision public metadata evidence is accepted only when explicitly pinned by CURATED_LICENSE_EVIDENCE.json under Security #150 Q3(B).')
[void]$notice.AppendLine('')
[void]$notice.AppendLine('==== Main project license ====')
[void]$notice.AppendLine((Get-Content -LiteralPath 'LICENSE' -Raw))
[void]$notice.AppendLine('')

$missing = @()
$manifest = @()
foreach ($dep in $deps) {
    $downloadJson = (& go mod download -json "$($dep.Path)@$($dep.Version)" 2>&1 | Out-String)
    if ($LASTEXITCODE -ne 0) { throw "go mod download failed for $($dep.Path)@$($dep.Version): $downloadJson" }
    $info = $downloadJson | ConvertFrom-Json
    if (-not $info.Dir) { throw "No module directory returned for $($dep.Path)@$($dep.Version)" }

    $licenseFiles = @(Get-ChildItem -LiteralPath $info.Dir -File | Where-Object {
        $_.Name -match '^(?i)(LICENSE|LICENCE|COPYING|NOTICE|COPYRIGHT)(\..*|-.*|_.*)?$'
    } | Sort-Object Name)
    if ($licenseFiles.Count -eq 0) {
        $licenseFiles = @(Get-ChildItem -LiteralPath $info.Dir -File -Recurse -Depth 2 | Where-Object {
            $_.Name -match '^(?i)(LICENSE|LICENCE|COPYING|NOTICE|COPYRIGHT)(\..*|-.*|_.*)?$'
        } | Sort-Object FullName | Select-Object -First 20)
    }

    $key = "$($dep.Path)@$($dep.Version)"
    if ($licenseFiles.Count -eq 0 -and $curated.ContainsKey($key)) {
        $evidence = $curated[$key]
        if ([string]::IsNullOrWhiteSpace([string]$evidence.spdx) -or [string]::IsNullOrWhiteSpace([string]$evidence.evidence_url)) {
            throw "Curated license evidence is incomplete for $key"
        }
        [void]$notice.AppendLine("==== $($dep.Path) $($dep.Version) ====")
        [void]$notice.AppendLine("SPDX license classification: $($evidence.spdx)")
        [void]$notice.AppendLine("Evidence type: $($evidence.evidence_type)")
        [void]$notice.AppendLine("Exact-revision evidence: $($evidence.evidence_url)")
        [void]$notice.AppendLine("Upstream license declaration: $($evidence.evidence_text)")
        if ($evidence.derivation) { [void]$notice.AppendLine("Derivation note: $($evidence.derivation)") }
        if ($evidence.canonical_terms_url) { [void]$notice.AppendLine("Canonical terms/provenance reference: $($evidence.canonical_terms_url)") }
        [void]$notice.AppendLine('')
        $manifest += [pscustomobject]@{
            module = $dep.Path
            version = $dep.Version
            go_sum = $dep.Sum
            license_files = @()
            license_expression = [string]$evidence.spdx
            evidence_type = [string]$evidence.evidence_type
            evidence_url = [string]$evidence.evidence_url
            evidence_text = [string]$evidence.evidence_text
        }
        continue
    }

    if ($licenseFiles.Count -eq 0) {
        $missing += $key
        continue
    }

    $safe = (($dep.Path + '@' + $dep.Version) -replace '[^A-Za-z0-9._@-]', '_')
    $dest = Join-Path $licensesDir $safe
    New-Item -ItemType Directory -Force -Path $dest | Out-Null

    [void]$notice.AppendLine("==== $($dep.Path) $($dep.Version) ====")
    $filesForManifest = @()
    foreach ($file in $licenseFiles) {
        $targetName = ($file.FullName.Substring($info.Dir.Length).TrimStart('\','/') -replace '[\\/]', '__')
        Copy-Item -LiteralPath $file.FullName -Destination (Join-Path $dest $targetName)
        $filesForManifest += $targetName
        [void]$notice.AppendLine("---- $targetName ----")
        try {
            [void]$notice.AppendLine((Get-Content -LiteralPath $file.FullName -Raw -ErrorAction Stop))
        } catch {
            [void]$notice.AppendLine('[Binary/non-text license file copied separately; see licenses directory.]')
        }
        [void]$notice.AppendLine('')
    }
    $manifest += [pscustomobject]@{
        module = $dep.Path
        version = $dep.Version
        go_sum = $dep.Sum
        license_files = $filesForManifest
        license_expression = $null
        evidence_type = 'module-license-file'
        evidence_url = $null
        evidence_text = $null
    }
}

$notice.ToString() | Set-Content -Encoding utf8 (Join-Path $OutDir 'THIRD_PARTY_NOTICES.txt')
$manifest | ConvertTo-Json -Depth 8 | Set-Content -Encoding utf8 (Join-Path $OutDir 'LICENSE_MANIFEST.json')

$missingPath = Join-Path $OutDir 'MISSING_LICENSE_EVIDENCE.txt'
if ($missing.Count -gt 0) {
    $missing | Set-Content -Encoding utf8 $missingPath
    'HOLD' | Set-Content -Encoding ascii (Join-Path $OutDir 'LICENSE_GATE_STATUS.txt')
} else {
    if (Test-Path -LiteralPath $missingPath) { Remove-Item -LiteralPath $missingPath -Force }
    'PASS' | Set-Content -Encoding ascii (Join-Path $OutDir 'LICENSE_GATE_STATUS.txt')
}

# GPL-family code is linked into the executable. Bundle the exact repository plus vendored
# Go dependency sources so the pre-release can accompany the executable with corresponding source.
# This evidence is generated even while the license gate is HOLD so Security can inspect the
# exact candidate without losing downstream SBOM/vulnerability evidence.
$stageRoot = Join-Path $env:RUNNER_TEMP 'dld051-corresponding-source'
if (Test-Path -LiteralPath $stageRoot) { Remove-Item -LiteralPath $stageRoot -Recurse -Force }
New-Item -ItemType Directory -Force -Path $stageRoot | Out-Null
$baseZip = Join-Path $env:RUNNER_TEMP 'dld051-source-base.zip'
if (Test-Path -LiteralPath $baseZip) { Remove-Item -LiteralPath $baseZip -Force }
& git archive --format=zip --output=$baseZip HEAD
if ($LASTEXITCODE -ne 0) { throw 'git archive failed' }
Expand-Archive -LiteralPath $baseZip -DestinationPath $stageRoot -Force
& go mod vendor
if ($LASTEXITCODE -ne 0) { throw 'go mod vendor failed' }
Copy-Item -LiteralPath 'vendor' -Destination (Join-Path $stageRoot 'vendor') -Recurse -Force
Copy-Item -LiteralPath $licensesDir -Destination (Join-Path $stageRoot 'EXACT_BINARY_LICENSES') -Recurse -Force
Copy-Item -LiteralPath (Join-Path $OutDir 'THIRD_PARTY_NOTICES.txt') -Destination (Join-Path $stageRoot 'THIRD_PARTY_NOTICES.txt') -Force
Copy-Item -LiteralPath (Join-Path $OutDir 'EXACT_BINARY_MODULES.tsv') -Destination (Join-Path $stageRoot 'EXACT_BINARY_MODULES.tsv') -Force
Copy-Item -LiteralPath (Join-Path $OutDir 'LICENSE_GATE_STATUS.txt') -Destination (Join-Path $stageRoot 'LICENSE_GATE_STATUS.txt') -Force
if (Test-Path -LiteralPath (Join-Path $OutDir 'CURATED_LICENSE_EVIDENCE.json')) {
    Copy-Item -LiteralPath (Join-Path $OutDir 'CURATED_LICENSE_EVIDENCE.json') -Destination (Join-Path $stageRoot 'CURATED_LICENSE_EVIDENCE.json') -Force
}
if (Test-Path -LiteralPath $missingPath) {
    Copy-Item -LiteralPath $missingPath -Destination (Join-Path $stageRoot 'MISSING_LICENSE_EVIDENCE.txt') -Force
}
$sourceZip = Join-Path $OutDir 'dld-051-corresponding-source.zip'
if (Test-Path -LiteralPath $sourceZip) { Remove-Item -LiteralPath $sourceZip -Force }
Compress-Archive -Path (Join-Path $stageRoot '*') -DestinationPath $sourceZip -CompressionLevel Optimal

$hash = Get-FileHash -LiteralPath $binary -Algorithm SHA256
$sourceHash = Get-FileHash -LiteralPath $sourceZip -Algorithm SHA256
@(
    "butler.exe`t$($hash.Hash.ToLowerInvariant())",
    "dld-051-corresponding-source.zip`t$($sourceHash.Hash.ToLowerInvariant())"
) | Set-Content -Encoding ascii (Join-Path $OutDir 'SHA256SUMS.txt')

Write-Host "Exact linked module count: $($deps.Count)"
Write-Host "Modules with acceptable license evidence: $($manifest.Count)"
Write-Host "Modules missing license evidence: $($missing.Count)"
if ($missing.Count -gt 0) {
    Write-Warning "License gate remains HOLD: $($missing -join ', ')"
} else {
    Write-Host 'License evidence complete for all exact linked modules.'
}
Write-Host "Corresponding-source bundle: $sourceZip"
