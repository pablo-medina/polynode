# polynode
Una herramienta de línea de comandos para administrar varias versiones de Node.

# Configuración
## Ubicación del espacio de trabajo
Esta herramienta utiliza un directorio como espacio de trabajo donde guarda la configuración, las versiones de Node descargadas, y la versión actual.
Por defecto, esta ubicación es c:\polynode, pero puede modificarse cambiando la variable de entorno POLYNODE_PATH.

## Configuración manual de %PATH%
Para el correcto funcionamiento de esta herramienta, debe configurarse manualmente la variable de entorno PATH, para que incluya el directorio %POLYNODE_PATH%\current.
Si por algún motivo no se puede acceder a la configuración de variables de entorno de Windows, se puede utilizar el comando ***setenv.cmd*** creado al ejecutar ```polynode init```.

# Comandos

| Comando                          | Descripción                                               |
| -------------------------------- | --------------------------------------------------------- |
| polynode init                    | Inicializa entorno (crea directorio y archivo setenv.cmd) |
| polynode install &lt;version&gt; | Instala la versión de Node indicada                       |
| polynode use &lt;version&gt;     | Cambia a la versión de Node indicada                      |
| polynode list                    | Lista las versiones de node disponibles localmente        |
| polynode version                 | Muestra la versión de Node utilizada actualmente          |

