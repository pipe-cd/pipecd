# System Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PipeCD Plugins Registry System                       │
└─────────────────────────────────────────────────────────────────────────┘

                           ┌──────────────────────┐
                           │  GitHub Releases     │
                           │  (API Endpoint)      │
                           └──────────┬───────────┘
                                      │
                ┌─────────────────────┘
                │
                ▼
    ┌────────────────────────────────┐
    │ GitHub Actions Workflow         │
    │ (Scheduled, On Release)         │
    │                                 │
    │ .github/workflows/              │
    │ update-plugins-registry.yaml    │
    └────────────┬───────────────────┘
                 │
                 ▼
    ┌────────────────────────────────┐
    │  Python Update Script           │
    │                                 │
    │  scripts/                       │
    │  update-plugins-registry.py     │
    │                                 │
    │  1. Fetch releases from GitHub  │
    │  2. Match tag patterns          │
    │  3. Extract versions            │
    │  4. Compare semantically        │
    │  5. Generate JSON & Markdown    │
    └────────────┬───────────────────┘
                 │
        ┌────────┴────────┐
        │                 │
        ▼                 ▼
    ┌──────────┐    ┌──────────────┐
    │ JSON API │    │ Markdown Doc │
    │          │    │              │
    │plugins.  │    │plugins.md    │
    │json      │    │              │
    └────┬─────┘    └────┬─────────┘
         │               │
         ├─────┬─────────┤
         │     │         │
         ▼     ▼         ▼
    ┌──────────────────────────────┐
    │  Validation & Testing         │
    │                              │
    │  scripts/                    │
    │  validate-plugins-registry.py│
    │  test_registry_scripts.py    │
    │  plugins.schema.json         │
    └──────────┬───────────────────┘
               │
        ┌──────┴──────┐
        │ Valid?      │
        └──────┬──────┘
               │
        ┌──────▼──────────────┐
        │ If Changed:         │
        │ - Commit & Push     │
        │ - Create PR         │
        └─────────────────────┘
               │
        ┌──────▼──────────────────────┐
        │  Distributed to Users        │
        │                             │
        │  • docs/plugins.md          │
        │  • docs/plugins.json        │
        │  • docs/PLUGINS_REGISTRY.md │
        │  • GitHub releases (cached) │
        └─────────────────────────────┘
               │
        ┌──────┴─────────────────┐
        │                        │
        ▼                        ▼
    ┌──────────────┐      ┌─────────────────┐
    │ End Users    │      │ Tools/Websites  │
    │              │      │                 │
    │ Visit docs/  │      │ Query JSON API  │
    │ plugins.md   │      │ Parse data      │
    │              │      │                 │
    │ Copy version │      │ Auto-detect     │
    │ from table   │      │ latest version  │
    └──────────────┘      └─────────────────┘
```

## Data Flow Diagram

```
                      Plugins in pipecd repo
                    (inline plugins with git tags)
                              │
                ┌─────────────┼─────────────┐
                │             │             │
                ▼             ▼             ▼
        kubernetes      terraform       wait
        (v0.1.0)        (v0.1.0)      (v0.1.0)
                │             │             │
                └─────────────┼─────────────┘
                              │
                              ▼
                  ┌──────────────────────┐
                  │  GitHub Releases API │
                  │  (query tags)        │
                  └──────────┬───────────┘
                             │
        ┌────────────────────┼────────────────────┐
        │                    │                    │
        ▼                    ▼                    ▼
    Match tag            Compare           Extract
    patterns           semantic ver.      version
    (e.g.,             (v0.1.0 >         string
    k8s/*)             v0.0.1)           from tag
        │                    │                    │
        └────────────────────┼────────────────────┘
                             │
                             ▼
                ┌──────────────────────────┐
                │ Latest Version Found     │
                │ e.g., v0.1.0             │
                └──────────┬───────────────┘
                           │
            ┌──────────────┼──────────────┐
            │              │              │
            ▼              ▼              ▼
        JSON             Markdown      Validation
        Generate         Generate      Schema Check
        plugins.json     plugins.md    plugins.schema.json
            │              │              │
            └──────────────┼──────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │ Commit and  │
                    │ Push to Git │
                    └─────────────┘
```

## Component Diagram

```
┌───────────────────────────────────────────────────────────────────────────┐
│                         Documentation Layer                               │
│                                                                           │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────────┐   │
│  │  QUICKSTART.md   │  │ PLUGINS_REGISTRY │  │ IMPLEMENTATION       │   │
│  │  (Quick ref)     │  │ .md (Full docs)  │  │ _SUMMARY.md (Details)│   │
│  └──────────────────┘  └──────────────────┘  └──────────────────────┘   │
└───────────────────────────────────────────────────────────────────────────┘
                                      │
┌───────────────────────────────────────────────────────────────────────────┐
│                         Registry Data Layer                               │
│                                                                           │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────────┐   │
│  │  plugins.json    │  │  plugins.md      │  │  plugins.schema.json │   │
│  │  (API)           │  │  (Human readable)│  │  (Validation)        │   │
│  └──────────────────┘  └──────────────────┘  └──────────────────────┘   │
└───────────────────────────────────────────────────────────────────────────┘
                                      │
┌───────────────────────────────────────────────────────────────────────────┐
│                         Processing Layer                                  │
│                                                                           │
│  ┌──────────────────────────┐    ┌─────────────────────────────────┐    │
│  │   Update Scripts         │    │    Validation & Testing         │    │
│  │                          │    │                                 │    │
│  │  • update-plugins-       │    │  • validate-plugins-registry    │    │
│  │    registry.py           │    │    .py                          │    │
│  │  • 350 lines             │    │  • 200 lines                    │    │
│  │  • Fetch releases        │    │  • test_registry_scripts.py     │    │
│  │  • Match patterns        │    │  • 300 lines                    │    │
│  │  • Compare versions      │    │                                 │    │
│  │  • Generate output       │    │                                 │    │
│  └──────────────────────────┘    └─────────────────────────────────┘    │
└───────────────────────────────────────────────────────────────────────────┘
                                      │
┌───────────────────────────────────────────────────────────────────────────┐
│                      Automation Layer (GitHub Actions)                    │
│                                                                           │
│  ┌──────────────────────────────────────────────────────────────────────┐ │
│  │  .github/workflows/update-plugins-registry.yaml                     │ │
│  │                                                                      │ │
│  │  Triggers:                                                          │ │
│  │  • On plugin release                                               │ │
│  │  • Every 6 hours (scheduled)                                       │ │
│  │  • On main release completion                                      │ │
│  │  • Manual trigger                                                  │ │
│  │                                                                      │ │
│  │  Process:                                                           │ │
│  │  • Checkout code                                                    │ │
│  │  • Setup Python                                                     │ │
│  │  • Run update script                                                │ │
│  │  • Check for changes                                                │ │
│  │  • If changed: Commit & Push                                        │ │
│  │  • Create PR for visibility                                         │ │
│  └──────────────────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────────────────┘
                                      │
┌───────────────────────────────────────────────────────────────────────────┐
│                         Integration Layer                                 │
│                                                                           │
│  ┌────────────────┐  ┌────────────────┐  ┌─────────────────────────┐    │
│  │ Make Targets   │  │ Direct Python  │  │ GitHub API Integration  │    │
│  │                │  │                │  │                         │    │
│  │ make gen/      │  │ python3 scripts│  │ Raw file access via     │    │
│  │ plugins-       │  │ /update-       │  │ raw.githubusercontent  │    │
│  │ registry       │  │ plugins-       │  │ .com                    │    │
│  │                │  │ registry.py    │  │                         │    │
│  │ make check/    │  │                │  │ CI/CD pipeline use      │    │
│  │ plugins-       │  │ python3 scripts│  │                         │    │
│  │ registry       │  │ /validate-     │  │                         │    │
│  │                │  │ plugins-       │  │                         │    │
│  │                │  │ registry.py    │  │                         │    │
│  └────────────────┘  └────────────────┘  └─────────────────────────┘    │
└───────────────────────────────────────────────────────────────────────────┘
                                      │
┌───────────────────────────────────────────────────────────────────────────┐
│                            End Users                                      │
│                                                                           │
│  • Browse docs/plugins.md                                                │
│  • Query docs/plugins.json API                                           │
│  • Download plugins from links                                           │
│  • Check GitHub releases page                                            │
└───────────────────────────────────────────────────────────────────────────┘
```

## Deployment Sequence

```
Time  Action                          Result
────  ──────────────────────────────  ──────────────────────────────
T+0   Plugin released on GitHub       Release tag created
      (e.g., pkg/app/.../v0.2.0)
      │
T+1   GitHub Actions triggered        Workflow starts
      (on:release)                    │
      │
T+5   Checkout & Setup Python         Environment ready
      │
T+10  Run update script               
      • Fetch GitHub releases
      • Match tag patterns
      • Extract v0.2.0 (new!)
      • Generate JSON & Markdown
      │
T+15  Validation                      Schema & semantic checks
      ✓ All checks pass
      │
T+20  Detect changes                  Version changed!
      • plugins.json updated
      • plugins.md updated
      │
T+25  Commit & Push                   
      git commit -m "chore: update plugins registry"
      git push origin master
      │
T+30  Workflow completes              SUCCESS
      │
      ✓ docs/plugins.json updated
      ✓ docs/plugins.md updated
      ✓ Changes committed to git
      ✓ Changes pushed to GitHub
      │
T+35  Users see updates
      • Visit docs/plugins.md: sees v0.2.0
      • API query: gets v0.2.0
      • Tools/websites: fetch updated data
```

## File Organization

```
pipecd/
├── docs/
│   ├── plugins.json                    ← Machine-readable registry
│   ├── plugins.md                      ← Human-readable doc
│   ├── plugins.schema.json             ← JSON schema
│   ├── PLUGINS_REGISTRY.md             ← Full documentation
│   └── [other existing files...]
│
├── scripts/
│   ├── update-plugins-registry.py      ← Main update script
│   ├── validate-plugins-registry.py    ← Validation script
│   ├── test_registry_scripts.py        ← Unit tests
│   ├── README.md                       ← Scripts documentation
│   └── [existing scripts...]
│
├── .github/
│   └── workflows/
│       ├── update-plugins-registry.yaml ← Automation workflow
│       └── [existing workflows...]
│
├── Makefile                            ← Updated with new targets
│
├── IMPLEMENTATION_SUMMARY.md           ← Implementation overview
├── PLUGINS_QUICKSTART.md               ← Quick start guide
├── IMPLEMENTATION_COMPLETE.md          ← Completion summary
│
└── [other existing files...]
```

## Technology Stack

```
┌─────────────────────────────────────┐
│         GitHub Actions              │
│  (Automation & Scheduling)          │
└────────────────┬────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────┐
│  Python 3.8+                        │
│  • requests library (API calls)     │
│  • json (data parsing)              │
│  • jsonschema (validation)          │
└────────────────┬────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────┐
│  GitHub API v3                      │
│  • Releases endpoint                │
│  • Tags endpoint                    │
│  • Rate limiting                    │
└────────────────┬────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────┐
│  Output Formats                     │
│  • JSON (machines)                  │
│  • Markdown (humans)                │
│  • JSON Schema (validation)         │
└─────────────────────────────────────┘
```

---

## Summary

This architecture provides:
- ✅ **Automated** plugin version tracking
- ✅ **Reliable** with comprehensive validation
- ✅ **Scalable** to any number of plugins
- ✅ **Documented** at every level
- ✅ **Tested** with unit and integration tests
- ✅ **Secure** with proper token handling
- ✅ **Integrated** with GitHub workflow

The system updates automatically with zero manual intervention, serving users through multiple channels (Markdown doc, JSON API, GitHub releases).
