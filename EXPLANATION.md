# Explicación del Proyecto

Este proyecto es un microservicio que centraliza la información sobre personas desaparecidas y encontradas en Venezuela. La idea es tener una única base de datos con toda la información, para que sea más fácil de consultar y de integrar con otras plataformas.

## ¿Qué hace el sistema?

El sistema permite registrar y buscar personas, ya sea que estén desaparecidas o que hayan sido encontradas. También permite registrar listas de personas en hospitales.

## ¿Cómo funciona?

El sistema tiene una API (una forma de que otros programas se comuniquen con él) que permite hacer lo siguiente:

- **Registrar una persona:** Permite agregar una nueva persona a la base de datos, con su información personal y su estado (desaparecido o encontrado).
- **Buscar una persona:** Permite buscar personas en la base de datos por su nombre, estado (desaparecido/encontrado) o de forma precisa utilizando su número de identificación nacional (Cédula de Identidad) a través de un parámetro de búsqueda (`?national_id=X`).
- **Registrar una lista de hospital:** Permite agregar una lista de personas que se encuentran en un hospital.

## Resolución de Conflictos de Ruteo (Net/HTTP)

Durante el desarrollo, detectamos un conflicto en el enrutador nativo de Go (`http.ServeMux`) entre las rutas `/persons/{id}/reports` (para ver reportes) y `/persons/national/{national_id}` (para buscar por cédula). Esto ocurría porque Go no podía discernir de forma única a qué ruta enviar una petición como `/persons/national/reports`.

Para solucionarlo de forma elegante y conforme a las mejores prácticas de arquitectura REST:
- Eliminamos la ruta conflictiva `/persons/national/{national_id}`.
- Unificamos la búsqueda en el endpoint principal de consulta: `GET /persons?national_id=X`. Ahora, si se proporciona la cédula como parámetro de consulta, el sistema devuelve directamente la lista con el registro que coincide con esa cédula (o vacío si no existe), evitando colisiones de ruteo y manteniendo una API limpia y estándar.

## Estructura de Datos (Modelos)

Hemos definido las entidades que representarán nuestra base de datos en el código:

- **Person (Persona):** Guarda el nombre, identificación nacional (cédula), fecha de nacimiento, género, último lugar y fecha vistos, y su estado actual (por ejemplo: 'missing', 'found', 'in_hospital').
- **Report (Reporte):** Permite almacenar reportes hechos por personas externas sobre un desaparecido o alguien encontrado, incluyendo el contacto de quien reporta.
- **Hospital (Hospital):** Contiene la información de contacto y ubicación de los centros de salud.
- **HospitalPerson (Personas en Hospital):** Vincula a las personas con el hospital donde están ingresadas, con detalles como fecha de ingreso, alta y notas médicas relevantes.

## Repositorios (Acceso a Datos)

Para interactuar con la base de datos, utilizamos un patrón de diseño llamado **Repositorio**. Esto nos permite separar la forma en que guardamos los datos de las reglas de negocio del sistema:

- **Definición de Contratos (Interfaces):** Creamos una lista de lo que nuestro sistema necesita hacer con los datos (por ejemplo: crear persona, buscar por ID, listar hospitales). Esto se encuentra en `domain/repository.go`.
- **Implementación Concreta (PostgreSQL):** Programamos la lógica real que se comunica con PostgreSQL. En Go, para evitar conflictos de nombres (ya que varias entidades necesitan operaciones con el mismo nombre, como `Create` o `GetByID`), separamos la lógica en tres estructuras distintas: `PostgresPersonRepository` (personas), `PostgresReportRepository` (reportes) y `PostgresHospitalRepository` (hospitales). Esto se encuentra en `adapters/postgres_repository.go`.

## ¿Qué tecnologías se usan?

- **Go:** Es el lenguaje de programación en el que está escrito el microservicio.
- **PostgreSQL:** Es la base de datos que se usa para guardar la información. La conexión con esta base de datos se realiza mediante un formato estándar llamado **Connection String** (URL de conexión), que agrupa todos los parámetros (usuario, contraseña, host, puerto, nombre de base de datos) en una sola cadena de texto (`postgres://...`). Esto facilita la configuración en entornos locales y servicios en la nube.
- **Swagger:** Es una herramienta que se usa para documentar la API, para que otros programadores puedan entender cómo usarla.

## Despliegue en Render

Para facilitar la publicación y despliegue del microservicio en internet usando **Render**, creamos un archivo llamado `render.yaml`. Este archivo es un Blueprint (plantilla) que le dice a Render exactamente cómo configurar, compilar e iniciar la aplicación.

Adicionalmente, para dar soporte **plug-and-play absoluto** (para que funcione de forma inmediata utilizando el comando predeterminado e inalterable de Render sin necesidad de cambiar nada manualmente en su interfaz), trasladamos el archivo principal de arranque de `cmd/main.go` a la raíz del repositorio (`./main.go`):
- **Resolución definitiva de error:** Al colocar `main.go` en la raíz del proyecto, el comando predeterminado de Render (`go build -tags netgo -ldflags '-s -w' -o app`) encuentra el código al instante, lo compila exitosamente, genera el binario `./app` y lo arranca sin fallar con errores de "archivos de Go no encontrados".
- **Variables de entorno:** Configura el puerto por defecto (`8080`) y define el espacio para colocar la conexión segura a la base de datos de PostgreSQL en producción (`DATABASE_URL`).
