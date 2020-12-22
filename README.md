# Rakoon back-end

  Rakoon est un bureau distant permettant de stocker et télécharrger des fichiers dans le cloud, sans limite de temps.
  Une gestion d'utilisateurs avec différents droits est implémentée et paramétable depuis un panneau de configuration.
  L'authentification se fait à l'aide de jsons web tokens.
  Ce dépot contient l'app back-end, en go.
  J'utilise gin-gonic et sqlx.

## Pré-requis
  
  Pour lancer l'application, il est nécessaire d'avoir:
  
  * Une base postgresql. (les script créeant les tables nécessaires se trouvent dans le dossier sql)

  * Les variables d'environnement suivantes de définies:
    `DB_PORT=5432`
    `DB_HOST=localhost`
    `DB_USER=tom`
    `DB_NAME=rakoon_db`
    `DB_PWD=qwerty`
    `SECRET_KEY=secretkey`
    
    

## Testing

  L'app est testée en end-to-end, on peut les lancer avec une base locale ou distante avec la commande suivante:

  * source .env-test; go test ./tests
