# Backend

Ce dépot à pour objectif de mettre sous format EPUB des "novels" actuellement disponible sur internet.
Les sites depuis lesquels il est actuellement possible de les récupérer sont :
- https://novelbin.me/
- https://readnovelfull.com/

## Utilisation

`docker-compose up` pour lancer la partie backend du projet, la partie frontend se trouve [ici](https://github.com/MrColorado/frontend).
Le docker compose lance actuellement plusieurs services :
- Une base de donnée : Postgresql
- Un S3 : Minio
- Un broqueur de message : Nats
- Un proxy pour permettre la conversion de gRPC : Envoy
- Un server permettant de servir le frontend
- Un Service permettant de récupérer les "novels" et les convertir sous format EPUB

## Commandes utilise pour le dev

Docker build :
- `docker build --build-arg DIRECTORY=server --file server/Dockerfile --progress=plain . -t server`
- `docker build --build-arg DIRECTORY=book-handler --file book-handler/Dockerfile --progress=plain . -t book-handler`

Récupérer le schema de la DB :
- `dc exec database pg_dump -U $POSTGRES_USER novel --schema-only`
- psql databasename < data_base_dump (c'est quoi déjà ?)
  
Générer les structure de sqlboiler :
- `sqlboiler -c ../schemas/sqlboiler.toml --wipe psql`

Se connecter sur le DB : 
- `dc exec database psql -h database -U $POSTGRES_USER novel`
- `docker exec -it database psql -h database -U $POSTGRES_USER novel`
- `psql --host=host.docker.internal --user=root_user --password novel`

pg_dump -U dbusername dbname > dbexport.pgsql
dc exec database pg_dump -U $POSTGRES_USER novel > dbexport.pgsql
dc exec database pg_dump -U $POSTGRES_USER postgres --schema-only


Il peut y avoir des conflits entre le docker psql et le service psql qui tourne sur wsl2 : 'sudo systemctl stop postgresql'