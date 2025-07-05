# Weather-Boy Implementation Plan

This document outlines the steps to complete the Weather-Boy backend, turning it into a fully functional weather data fusion and risk analysis service.

## Project Goals

1.  **Data Ingestion:** Fetch, parse, and store data from multiple IMD sources: Bulletin, Radar, Nowcast, District Warnings, and River Basin Forecasts.
2.  **Risk Fusion:** Calculate a comprehensive, location-aware risk score by fusing the different data sources.
3.  **API Exposure:** Expose both the raw/parsed data and the final risk score via a clean, versioned API. The risk score response **must** include a breakdown of its components.

## Current Status (As of 2025-07-04)

-   **Nowcast:** Fetching and storing successfully.
-   **District Warnings:** Fetching and storing successfully.
-   **Bulletin PDF:** Downloading the raw PDF, but not parsing it.
-   **Radar:** Not implemented.
-   **River Basin:** Not implemented.
-   **Risk Score:** Implemented but only using Nowcast data. It is incomplete.

---

## ✅ TODO

### Phase 1: Radar Data Pipeline

-   **[DONE]** Create `internal/fetch/radar.go` to download the latest Doppler radar PNG from `mausam.imd.gov.in`.
-   **[DONE]** Research and document the IMD's color-to-dBZ mapping and the radar's pixel-to-kilometer scale.
-   **[DONE]** Create `internal/parse/radar.go` to analyze the PNG image, find the max dBZ value within a 40km radius of the target location, and extract the timestamp.
-   **[DONE]** Create a database migration and repository for a `radar` table to store the parsed data (`location`, `max_dbz`, `captured_at`).


### Phase 2: Bulletin Parsing Pipeline (AI-Powered)

-   **[DONE]** Add a Go library for PDF-to-text extraction.
-   **[DONE]** Add a Go library compatible with the OpenAI/Gemini API.
-   **[DONE]** Create `internal/parse/bulletin.go` which will:
    1.  Extract text from the downloaded PDF.
    2.  Call the Gemini Flash model with a prompt to get a structured forecast for the target city.
-   **[DONE]** Create a database migration and repository for a `bulletin_parsed` table to store the AI-generated summary.
-   **[DONE]** Add `OPENAI_API_KEY` to `.env.example` and load it in the config.


### Phase 3: River Basin Data Pipeline

-   **[DONE]** Create `internal/fetch/riverbasin.go` to fetch data from the `basin_qpf_api.php`.
-   **[DONE]** Create a database migration and repository for a `river_basin_qpf` table.


### Phase 4: Integration & Finalization

-   **[DONE]** Update `internal/score/score.go` to incorporate the new radar, bulletin, and river basin data into the risk calculation.
-   **[DONE]** Update the `/v1/risk/:loc` handler to return a `breakdown` field detailing how each component contributed to the final score.
-   **[DONE]** Implement the `/v1/bulletin/:loc` and `/v1/radar/:loc` handlers to serve the newly parsed data.
-   **[DONE]** Implement a new `/v1/riverbasin/:loc` handler.
-   **[IN PROGRESS]** Update the main scheduler in `internal/scheduler/jobs.go` to run all the new fetch and parse jobs.

### Phase 5: AWS/ARG Data Pipeline

-   **[DONE]** Create `internal/fetch/awsarg.go` to fetch data from the `aws_data_api.php`.
-   **[DONE]** Create a database migration and repository for a `aws_arg` table.
-   **[DONE]** Update `internal/score/score.go` to incorporate the new AWS/ARG data into the risk calculation.
-   **[DONE]** Implement a new `/v1/awsarg/:loc` handler.
-   **[DONE]** Update the main scheduler in `internal/scheduler/jobs.go` to run the new fetch job.

---

## 🔄 In Progress

-   **[IN PROGRESS]** Finalize scheduler updates to ensure all new jobs (radar, river basin, AWS/ARG) are registered and monitored.

---

## ✔️ Done

-   Radar data pipeline
-   Bulletin AI parsing pipeline
-   River Basin QPF pipeline
-   AWS/ARG ingestion pipeline
-   Risk score integration with all new data sources
-   New API endpoints for bulletin, radar, river basin, and AWS/ARG
