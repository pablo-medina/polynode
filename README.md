# polynode
Una herramienta de línea de comandos para administrar varias versiones de Node.

# Configuración
## Ubicación del espacio de trabajo
Esta herramienta utiliza un directorio como espacio de trabajo donde guarda la configuración, las versiones de Node descargadas, y la versión actual.
Por defecto, esta ubicación es c:\polynode, pero puede modificarse cambiando la variable de entorno POLYNODE_PATH.

## Inicialización del espacio de trabajo
Utilizar el comando ```poly init``` para que se inicialice el espacio de trabajo con el repositorio de versiones vacío.

## Configuración manual de %PATH%
Para el correcto funcionamiento de esta herramienta, debe configurarse manualmente la variable de entorno PATH, para que incluya el directorio %POLYNODE_PATH%\current.
El comando ```poly check``` ayuda a verificar si dicho PATH está correctamente configurado..

## Proxy
En caso de necesitar indicar la configuración del proxy, se debe crear el archivo ***proxy.json** en el espacio de trabajo (por ejemplo: c:\polynode\proxy.json) con esta estructura:

```json
{
    "http_proxy": "http://<url>:<puerto>"
}
```
De esta manera, se va tomar la URL de proxy indicada en el archivo antes mencionado.

# Comandos

| Comando                      | Descripción                                                         |
| ---------------------------- | ------------------------------------------------------------------- |
| poly install &lt;version&gt; | Instala la versión de Node indicada                                 |
| poly use &lt;version&gt;     | Cambia a la versión de Node indicada                                |
| poly list                    | Lista las versiones de node disponibles localmente                  |
| poly version                 | Muestra la versión de Node utilizada actualmente                    |
| poly uninstall               | Desinstala la versión de Node indicada del repositorio local        |
| poly proxy <url>             | Definir la URL del proxy                                            |
| poly check                   | Verifica la instalación de polynode                                 |
| poly backup                  | Realiza una copia de seguridad de la instalación actual de polynode |
| poly shell                   | Abre un shell con la versión actual de Node.js configurada en el PATH |
| poly help                    | Mostrar ayuda de línea de comandos                                  |
