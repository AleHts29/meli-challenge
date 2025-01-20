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
│   └── server/
│       ├── server/
│       │   └──  handler.go  # Manejo de peticiones HTTP relacionadas con IPs
│       └── main.go          # Punto de entrada de la aplicación
├── internal/
│   ├── ipinfo/
│   │   ├── service.go       # Lógica de negocio para resolver información de IPs
│   │   ├── repository.go    # Interacción con APIs externas (países y monedas)
│   │   └── blocklist.go     # Lógica para manipular la lista de bloqueos
│   ├── models/
│   │   └── model.go         # Modelos estructuras
│   └── config/
│       └── config.go        # Configuración de la aplicación (API Keys, etc.)
├── pkg/
│   ├── api/
│   │   ├── countries.go     # Cliente para la API de países
│   │   └── currencies.go    # Cliente para la API de cotización de monedas
│   └── cache/
│       └── cache.go         # Implementación de caché en memoria
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
   
5. Agregar el archivo `IP2LOCATION-LITE-DB1.BIN` en la ruta raiz del proyecto.
   
6. Inicia el servidor:
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

## DB de IPs

Con fin de trabajar con un listado de IPs para los diferentes paises en los que opera **MercadoLibre**, se esta usando una base de datos en formato binario `IP2LOCATION-LITE-DB1.BIN` que se descargo de IP2Location (recurso es gratuito), tambien se agrega un archivo .CSV de `LACNIC` para buscar una IP de prueba y poder realizar las requests correspondientes.

Puedes consultar más detalles en [LANIC](https://ftp.lacnic.net/pub/stats/lacnic/) y [IP2_LOCATION](https://lite.ip2location.com/database-download)




## Ejemplo Lista de IPs de pruebas
```csv
lacnic|BR|ipv4|45.70.232.0|1024|20170913|allocated|127887
lacnic|EC|ipv4|45.70.236.0|1024|20170913|allocated|279839
lacnic|BR|ipv4|45.70.244.0|1024|20170915|allocated|278989
lacnic|BR|ipv4|45.70.248.0|1024|20170925|allocated|220610
lacnic|BR|ipv4|45.70.252.0|1024|20170925|allocated|239922
lacnic|EC|ipv4|45.71.0.0|1024|20171004|allocated|279375
lacnic|BR|ipv4|45.71.4.0|256|20170914|assigned|103789
lacnic|AR|ipv4|45.71.5.0|256|20171123|assigned|84411
lacnic|BR|ipv4|45.71.6.0|256|20170928|assigned|267952
lacnic|CO|ipv4|45.71.7.0|256|20171124|assigned|280200
lacnic|CL|ipv4|45.71.8.0|1024|20170915|assigned|259232
lacnic|BR|ipv4|45.71.12.0|1024|20170914|allocated|276567
```