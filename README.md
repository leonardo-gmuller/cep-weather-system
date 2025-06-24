# 📦 Project: CEP Weather System

A distributed system that, given a Brazilian **CEP (postal code)**, returns the **city name** and the **current temperature** using two independent microservices:

- **Service A**: Validates the CEP and retrieves the city.
- **Service B**: Retrieves the temperature based on the city.

The project follows **Clean Architecture**, uses HTTP communication between services, and features **distributed tracing with OpenTelemetry + Zipkin**.

---

## 🧱 Architecture

```text
[User/Postman]
      |
      v
[Service A: /service-a/v1/validate-cep] ---> [Service B: /weather/{cep}]
```

- Distributed tracing using OpenTelemetry
- Traces are visualized with **Zipkin**
- Services are written in Go, with Docker used for observability tools

---

## 🚀 Getting Started

### Prerequisites

- Docker and Docker Compose installed
- Go 1.23+ installed locally

---

### Environment Configuration

Make sure to create a `.env` file at the root of your project with the necessary environment variables used by both services.

Use `.env-template`

> 🔧 Update the values as needed based on whether you're running inside or outside Docker.

---

### Run only Zipkin and OpenTelemetry Collector

```bash
docker-compose up -d otel-collector zipkin
```

> 💡 Use this command if you're running services A and B manually on your machine (outside Docker).

---

### Run everything with Docker (Services A & B + Otel + Zipkin)

```bash
docker-compose up --build
```

- **Service A**: http://localhost:8080  
- **Service B**: http://localhost:8081  
- **Zipkin UI**: http://localhost:9411  

---

## 🌐 API Endpoints

### Service A

- `POST /service-a/v1/validate-cep`

**Request body:**

```json
{
  "cep": "01001000"
}
```

**Sample response:**

```json
{
  "city": "São Paulo",
  "temp_c": 22.5,
  "temp_f": 72.5,
  "temp_k": 295.65
}
```

---

## 📡 Observability

- OpenTelemetry is configured with OTLP over gRPC
- Zipkin is used to visualize all trace data
- The trace context is propagated across services

👉 Access Zipkin at: http://localhost:9411

---

## 🛠️ Technologies Used

- Go
- OpenTelemetry SDK for Go
- OTLP gRPC Protocol
- Zipkin
- Docker / Docker Compose
- Chi Router

---

## 📁 Project Structure

```bash
.
├── service-a/
│   ├── cmd/
│   └── internal/
├── service-b/
│   ├── cmd/
│   └── internal/
├── otel-collector-config.yaml
├── docker-compose.yml
└── README.md
```

---
