## 1. what is weather-boy?

**weather-boy** is a two-part monorepo that:

1. **pulls multiple IMD data sources** (district bulletin PDF, doppler radar PNG, nowcast JSON, etc.),
2. **surf-faces the raw numbers / phrases** so humans can see exactly what the feeds say, and
3. **distills them into a single traffic-light “risk” verdict** (green / yellow / orange / red) for the next 0–4 h and for the rest of the day.

the frontend shows both the *raw* feeds (latest radar image, bulletin excerpt, POP graph, etc.) **and** the fused verdict badge.

### outputs

| endpoint                 | ttl / cadence | shows                                                                 |
|--------------------------|---------------|-----------------------------------------------------------------------|
| `/v1/risk/:loc`          | 30 s cache    | `level, score, breakdown` — the traffic-light verdict + subscores     |
| `/v1/bulletin/:loc`      | 24 h          | parsed page from bulletin PDF (phrase & confidence)                   |
| `/v1/radar/:loc`         | 5 min         | `{max_dbz, bearing, range_km, timestamp}`                             |
| `/v1/nowcast/:loc`       | 15 min        | array of 0-4 h `{lead_min, pop, mm_per_hr}` buckets                   |

---

## 2. repo layout (top level)

weather-boy/
├── AGENTS.md            ← this file
├── fe/                  ← vite + react + tailwind + shadcn/ui
└── be/                  ← go 1.22 + fiber api

---

## 3. data sources

| source                              | url template                                                | cadence |
|-------------------------------------|-------------------------------------------------------------|---------|
| **bulletin pdf**                    | `https://…/mcdata/{PdfSlug}`                                | daily ~18:30 IST |
| **doppler radar**                   | `https://mausam.imd.gov.in/{RadarCode}_latest_ref.png`      | 5 min |
| **metnet nowcast** (fallback)       | `https://nowcast-ml-api.metnet.cloud/{lat},{lon}`           | 15 min |
| **districtWarning / nowcast json**  | IMD API                                                     | hourly / 30 min |

### hard-coded cities (v0)

vadodara | 22.30 N, 73.20 E  (primary)
mumbai   | 19.08 N, 72.88 E
thane    | 19.22 N, 72.97 E
pune     | 18.52 N, 73.85 E

add a city by appending to `internal/config/locations.go`.

---

## 4. risk scoring v0

score = 0
if bulletin says heavy/very-heavy       → +0.4
if radar max_dBZ ≥45 within 40 km       → +0.4
if nowcast POP₁ₕ ≥0.7                  → +0.2

threshold → traffic-light:

| score | level  |
|-------|--------|
| ≥0.8  | RED    |
| ≥0.5  | ORANGE |
| ≥0.3  | YELLOW |
| else  | GREEN  |

return `score` and per-sensor breakdown so UI can justify itself.

---

## 5. agent etiquette

1. **no secrets in browser** — openai key only on backend.
2. **one feature, one branch** — conventional commits (`feat: …`, `fix: …`).
3. **tests before tokens** — unit-test parse funcs; llm calls stubbed behind an interface.
4. **cron jitter** — every fetch job sleeps ±30 s to avoid hammering IMD.
5. **schema via golang-migrate** — never hand-edit DB.
6. **style consistency** — backend logs use zap structured, frontend uses shadcn cards/badges.

---

## 6. open tasks (2025-06-18)

- [ ] implement `internal/fetch/bulletin.go` (pdf grab + cache)
- [ ] implement `internal/parse/bulletin.go` (slice + o3)
- [ ] implement `internal/fetch/radar.go` & parser
- [ ] wire cron jobs in `internal/scheduler/jobs.go`
- [ ] expose new `/v1/risk/:loc` route
- [ ] seed frontend with mock JSON until IMD API lives

after IMD whitelist:

- [ ] `fetch/district_warning.go`
- [ ] swap MetNet for official nowcast JSON

---

> questions? open an issue or ping ayaan.
