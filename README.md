# This workspace provides 2 utilities:
- [Extractor](extractor/README.md)
- [Loader](loader/README.md)

This project is intended to be used for anyone that requires to extract/load massive amounts of data from/into MongoDb in batch processes.

This is helpful on environments such as:
- Airflow
- Dagster
- Prefect

Any pipeline orchestrator can benefit from it, those are generally written in python and extracting/loading data using python code is convenient but not performant.
This package have dockerized version too, this ensures you benefit from this by using a smaller container with faster perfomance.
