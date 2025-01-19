# Meli Challenge

Este proyecto tiene como objetivo coordinar acciones de respuesta ante posibles fraudes en transacciones de compra, búsqueda y pago en la plataforma. Para ello, se desarrolló una API REST que proporciona información contextual sobre una dirección IP, permite gestionar bloqueos y notifica a equipos interesados de manera descentralizada.

## Funcionalidades

1. **Consulta de información por dirección IP**:
    - Determinar el país asociado a una dirección IP.
    - Mostrar información del país, incluyendo la cotización en dólares de su moneda.

2. **Bloqueo de IPs**:
    - Registrar una dirección IP en una lista de bloqueos, impidiendo nuevas consultas de información asociada.

3. **Notificaciones descentralizadas**:
    - Informar sobre bloqueos de IP a equipos interesados de forma descentralizada, sin conocer de antemano los destinatarios.

## APIs Externas Utilizadas

Para implementar estas funcionalidades, se consumen las siguientes APIs externas:
- **API de Países**: Información sobre la geolocalización de direcciones IP.
- **API de Cotización de Monedas**: Información sobre las tasas de cambio de divisas en relación con el dólar.

Puedes consultar más detalles en [APIs de MELI](https://developers.mercadolibre.com.ar/es_ar/ubicacion-y-monedas#close).

---

## Requisitos del Proyecto

- Evitar llamadas innecesarias a APIs externas mediante el uso de caché.
- Garantizar la persistencia de datos (como la lista de bloqueos) ante un reinicio de la aplicación.
- Implementar una arquitectura escalable y modular basada en microservicios.
- Enviar notificaciones de forma descentralizada, asegurando independencia entre los servicios involucrados.

---

## Estructura del Proyecto

El proyecto está organizado en tres directorios principales: `cmd/`, `internal/` y `pkg/`. A continuación, se detalla la estructura del proyecto:

```text
mercado-libre-api/
├── cmd/
│   └── api/
│       ├── main.go          # Punto de entrada de la aplicación
├── internal/
│   ├── ipinfo/
│   │   ├── handler.go       # Manejo de peticiones HTTP relacionadas con IPs
│   │   ├── service.go       # Lógica de negocio para resolver información de IPs
│   │   ├── repository.go    # Interacción con APIs externas (países y monedas)
│   │   └── model.go         # Estructuras de datos relacionadas con IPs y países
│   ├── blocklist/
│   │   ├── handler.go       # Gestión de bloqueos de IP
│   │   ├── service.go       # Lógica para manipular la lista de bloqueos
│   │   ├── repository.go    # Persistencia de la lista de bloqueos
│   │   └── model.go         # Estructuras relacionadas con IPs bloqueadas
│   ├── notifications/
│   │   ├── service.go       # Lógica para envío de notificaciones descentralizadas
│   │   └── model.go         # Estructuras para gestionar las notificaciones
│   └── config/
│       └── config.go        # Configuración de la aplicación (API Keys, etc.)
├── pkg/
│   ├── api/
│   │   ├── countries.go     # Cliente para la API de países
│   │   └── exchange.go      # Cliente para la API de cotización de monedas
│   ├── cache/
│   │   └── cache.go         # Implementación de caché en memoria
│   ├── errors/
│   │   └── errors.go        # Manejo centralizado de errores
│   └── logger/
│       └── logger.go        # Logging centralizado para toda la aplicación
├── go.mod                   # Configuración del módulo de Go
└── go.sum                   # Archivo de dependencias
```

---

## Configuración

1. Clona este repositorio:
   ```bash
   git clone https://github.com/AleHts29/meli-challenge.git
   ```
2. Ve al directorio del proyecto:
   ```bash
   cd meli-challenge
   ```
3. Instala las dependencias:
   ```bash
   go mod tidy
   ```
4. Configura las variables de entorno en el archivo `config/config.go` para incluir:
    - API Keys necesarias para acceder a las APIs externas.
    - Configuración del puerto del servidor.

5. Inicia el servidor:
   ```bash
   go run cmd/api/main.go
   ```

---

## Tecnologías Utilizadas

- **Lenguaje**: Go (Golang)
- **Framework**: Gin para la creación de APIs REST.
- **Caché**: Implementación propia en memoria para reducir el tráfico hacia APIs externas.
- **Logs**: Módulo personalizado de logging centralizado.
- **Persistencia**: Archivos locales o bases de datos (según la implementación).

---

## Mejoras Futuras

- Integrar un sistema de mensajería como Kafka o RabbitMQ para manejar notificaciones de manera más eficiente.
- Incorporar autenticación y autorización para proteger la API.
- Implementar métricas y monitorización para analizar el rendimiento.