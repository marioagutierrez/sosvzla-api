# Microservicio de Búsqueda de Personas (SOS Venezuela)

Este microservicio en Go centraliza la información sobre personas desaparecidas y encontradas, así como listados de personas en hospitales, con el fin de unificar la información de distintas plataformas y facilitar su integración pública. El sistema cuenta con documentación interactiva en **Swagger** y utiliza una base de datos **PostgreSQL**.

## Características Principales

- **Registro de Personas:** Almacenamiento unificado de personas desaparecidas, encontradas o ingresadas en centros de salud.
- **Reportes Ciudadanos:** Registro de reportes independientes realizados por terceros sobre el paradero o avistamiento de personas.
- **Gestión de Hospitales:** Registro de centros de salud y listado de personas ingresadas en ellos para simplificar la búsqueda médica.
- **Arquitectura Limpia:** Organización desacoplada mediante el patrón de Repositorio y Servicios, facilitando extensiones y pruebas.
- **API RESTful e Interactiva:** Documentación integrada y probada a través de Swagger.

## Requisitos Previos

Asegúrate de contar con los siguientes componentes en tu sistema:

- **Go** (Versión 1.22 o superior recomendado)
- **PostgreSQL** (Base de datos activa)

## Estructura de la Base de Datos

Antes de iniciar el servicio, debes inicializar la estructura de la base de datos PostgreSQL utilizando el script SQL provisto en la carpeta del proyecto:

```bash
# Ejecuta el script SQL en tu cliente de base de datos preferido o consola de PostgreSQL
psql -u tu_usuario -d tu_base_de_datos -f build/database.sql
```

## Configuración

1. Copia el archivo de configuración de ejemplo para crear tu archivo `.env`:
   ```bash
   cp .env.example .env
   ```
2. Abre el archivo `.env` recién creado y reemplaza la cadena de conexión por la de tu base de datos real:
   ```env
   DATABASE_URL=postgres://usuario:contraseña@localhost:5432/nombre_db?sslmode=disable
   ```

## Cómo Iniciar el Proyecto

1. Instala todas las dependencias requeridas del módulo:
   ```bash
   go mod tidy
   ```
2. Compila el microservicio para comprobar que todo esté correcto:
   ```bash
   go build -o build/main main.go
   ```
3. Ejecuta el servidor de desarrollo:
   ```bash
   go run main.go
   ```

El servidor iniciará por defecto en el puerto `8080`.

---

## Documentación de la API (Swagger)

Cuando el servicio esté en ejecución, puedes acceder a la interfaz interactiva de Swagger UI para probar todos los endpoints disponibles:

👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

### Resumen de Endpoints Disponibles

#### Personas (`/persons`)
- `POST /persons` - Registrar una persona (desaparecida, encontrada, etc.).
- `GET /persons/{id}` - Obtener detalles de una persona por su ID de base de datos.
- `PUT /persons/{id}` - Actualizar detalles o estado de una persona.
- `DELETE /persons/{id}` - Eliminar un registro de persona.
- `GET /persons` - Listar personas (con paginación y filtro opcional por `status`).
- `GET /persons?national_id=X` - Buscar una persona de forma precisa por su Cédula de Identidad.

#### Reportes (`/reports`)
- `POST /reports` - Registrar un reporte de avistamiento o información de búsqueda.
- `GET /persons/{id}/reports` - Listar todos los reportes asociados a una persona específica.

#### Hospitales (`/hospitals`)
- `POST /hospitals` - Registrar un centro de salud.
- `GET /hospitals` - Listar todos los centros de salud registrados.
- `POST /hospitals/{id}/admit` - Registrar el ingreso de una persona a un hospital específico (actualiza automáticamente el estado de la persona a `in_hospital`).
- `GET /hospitals/{id}/persons` - Listar todas las personas ingresadas en un hospital específico.

---

## Regenerar Documentación de Swagger

Si realizas modificaciones en los endpoints o estructuras de datos y deseas actualizar la especificación de Swagger, puedes regenerarla utilizando la herramienta de Swag:

```bash
go run github.com/swaggo/swag/cmd/swag init -g main.go
```

## Autor

Desarrollado y mantenido por:
- **Mario Gutiérrez** - [@marioagutierrez](https://github.com/marioagutierrez)

## Licencia

Este proyecto está bajo la Licencia MIT. Para más detalles, consulta el archivo de licencia correspondiente.
