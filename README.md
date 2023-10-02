# This workspace provides 2 utilities:
- [Extractor](https://github.com/farovictor/MongoDbExtractor)
- [Loader](https://github.com/farovictor/MongoDbLoader)

This project is intended to be used for anyone that requires to extract/load massive amounts of data from/into MongoDb in batch processes.

## Requirements
This project makes use of [goenv](https://github.com/go-nv/goenv), a go version management tool to use the proper go version.

- go 1.20.0

## Container
To build images you can either run docker-compose file or use the dockerfiles in subpackages from repo root context. Extract and Loader requires the workspace to be entirely copied so they can access the driver package.
p.s.: At this commit the packages are not officially published, thats why the docker context requires the workspace directory to be fully copied to container.

## Environments
Any of the most used workflow/pipeline orchestrator out there in the market can benefit from Extractor/Loader package, those are generally written in python, which is convenient but not really performant.

Most used workflow orchestrators in the market:
- [Airflow](https://airflow.apache.org/)
- [Dagster](https://dagster.io/)
- [Prefect](https://www.prefect.io/)
- [Mage](https://docs.mage.ai/)

That is when those packages become really handy. You may use them as packages on your go projects, build and use as binary or just use images with already built binaries in your kubernetes pods.
This packages are simple, compact and make use of go power and concurrency to delivery the fastes and most performant routines.
