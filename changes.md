# Registro de Cambios

## 25-06-2026

- **CREAR:** Creación de la estructura inicial del proyecto.
  - **Propósito:** Establecer una base organizada para el microservicio, siguiendo las convenciones de Go.
  - **Lógica:** Se crearon los directorios `cmd`, `internal`, `adapters`, `app`, `domain`, `docs` y `build` para separar las responsabilidades del código.
- **CREAR:** Inicialización del módulo de Go.
  - **Propósito:** Permitir la gestión de dependencias del proyecto.
  - **Lógica:** Se ejecutó `go mod init` para crear el archivo `go.mod`.
- **CREAR:** Creación del archivo `main.go`.
    - **Propósito:** Establecer el punto de entrada de la aplicación.
    - **Lógica:** Se creó un archivo `main.go` simple en el directorio `cmd`.
- **CREAR:** Creación del archivo `README.md`.
    - **Propósito:** Proporcionar una descripción inicial del proyecto y cómo empezar a utilizarlo.
    - **Lógica:** Se creó un archivo `README.md` con secciones para la descripción, requisitos, instalación y documentación de la API.
- **CREAR:** Creación de los archivos de dominio `person.go`, `report.go` y `hospital.go`.
    - **Propósito:** Definir los modelos de datos que se usarán en el microservicio.
    - **Lógica:** Se definieron las estructuras `Person`, `Report`, `Hospital` y `HospitalPerson` en el paquete `domain`.
- **MODIFICAR:** Cambiada la conexión a la base de datos para usar un connection string (`DATABASE_URL`).
    - **Propósito:** Simplificar la configuración de base de datos facilitando la integración con servicios en la nube (ej. Heroku, Render, AWS RDS) que típicamente proveen un URL de conexión completo.
    - **Lógica:** Se modificó `internal/config/config.go` para cargar la variable `DATABASE_URL` únicamente, y se actualizó `.env.example` para reflejar el formato `postgres://usuario:contraseña@host:puerto/dbname?sslmode=disable`.
- **CREAR:** Creación de los archivos `.env.example`, `.gitignore` y `database.sql`.
    - **Propósito:** Facilitar la configuración local y asegurar que los archivos sensibles o innecesarios no se suban al repositorio.
    - **Lógica:** Se creó un archivo de configuración de ejemplo, un `.gitignore` y el esquema SQL inicial de la base de datos.
- **CREAR:** Creación de interfaces de repositorio `domain/repository.go`.
    - **Propósito:** Definir los métodos necesarios para la manipulación de datos de personas, reportes y hospitales de forma desacoplada.
    - **Lógica:** Se declararon las interfaces `PersonRepository`, `ReportRepository` y `HospitalRepository`.
- **CREAR:** Implementación de repositorios en Postgres en `adapters/postgres_repository.go`.
    - **Propósito:** Implementar la lógica de base de datos de manera desacoplada para evitar colisiones de firmas de métodos en Go.
    - **Lógica:** Se crearon tres estructuras separadas: `PostgresPersonRepository`, `PostgresReportRepository` y `PostgresHospitalRepository` para cumplir con las interfaces de dominio correspondientes de forma limpia y tipada, eliminando colisiones de nombres como `Create` y `GetByID`.
- **RESOLVER:** Corrección de conflicto de rutas en `http.ServeMux` de Go.
    - **Propósito:** Evitar el pánico al iniciar el servidor debido a rutas superpuestas que no se podían discernir de forma única (`GET /persons/{id}/reports` vs `GET /persons/national/{national_id}`).
    - **Lógica:** Se eliminó la ruta `/persons/national/{national_id}` de `adapters/http_handler.go` y en su lugar se integró la búsqueda por cédula/ID nacional como un parámetro de consulta (query parameter) en la ruta existente `GET /persons?national_id=X`. Esto es más elegante, RESTful y elimina por completo el conflicto de ruteo.
- **MODIFICAR:** Traducción y actualización de `README.md`.
    - **Propósito:** Ofrecer la documentación completa de inicio rápido, configuración, detalles de API y endpoints en español, y configurar el perfil de autor correcto en GitHub.
    - **Lógica:** Se reescribió `README.md` enteramente en español y se vinculó la cuenta de GitHub de Mario Gutiérrez (`https://github.com/marioagutierrez`) como autor principal del repositorio.
- **CREAR:** Archivo de configuración para despliegues en Render `render.yaml`.
    - **Propósito:** Corregir el fallo de compilación en Render (`no Go files in ...`) causado por la ausencia de archivos Go en la raíz del repositorio, especificando de forma explícita la ruta al punto de entrada.
    - **Lógica:** Se configuró un Blueprint de Render en `render.yaml` especificando el comando de compilación explícito `go build -tags netgo -ldflags '-s -w' -o app cmd/main.go` y el comando de inicio `./app` con las variables de entorno asociadas.

