# 📊 ServidorDeReproducciones - API REST Simple

## ✅ Arquitectura Correcta

Este servidor **SOLO** se encarga de:
1. Almacenar reproducciones de canciones
2. Consultar reproducciones por usuario

**NO contiene:**
- ❌ Cliente
- ❌ Archivos de audio
- ❌ Lógica de streaming
- ❌ Comunicación con otros servidores

---

## 🚀 Ejecución

```bash
cd ServidorDeReproducciones_Correcto/main
go run servidor.go
```

El servidor se inicia en: **http://localhost:3000**

---

## 📡 Endpoints

### 1. POST /reproducciones
**Función:** Almacenar una nueva reproducción

**Body (JSON):**
```json
{
  "idUsuario": 1,
  "idCancion": 1,
  "titulo": "Believer",
  "artista": "Imagine Dragons",
  "genero": "Rock",
  "idioma": "Inglés"
}
```

**Respuesta:**
```json
{
  "id": 3,
  "idUsuario": 1,
  "idCancion": 1,
  "titulo": "Believer",
  "artista": "Imagine Dragons",
  "genero": "Rock",
  "idioma": "Inglés",
  "fechaHora": "2025-11-01T15:30:00Z"
}
```

**Ejemplo con curl:**
```bash
curl -X POST http://localhost:3000/reproducciones \
  -H "Content-Type: application/json" \
  -d '{
    "idUsuario": 1,
    "idCancion": 1,
    "titulo": "Believer",
    "artista": "Imagine Dragons",
    "genero": "Rock",
    "idioma": "Inglés"
  }'
```

---

### 2. GET /reproducciones
**Función:** Obtener todas las reproducciones

**Respuesta:**
```json
[
  {
    "id": 1,
    "idUsuario": 1,
    "idCancion": 1,
    "titulo": "Believer",
    "artista": "Imagine Dragons",
    "genero": "Rock",
    "idioma": "Inglés",
    "fechaHora": "2025-11-01T10:30:00Z"
  },
  ...
]
```

**Ejemplo:**
```bash
curl http://localhost:3000/reproducciones
```

---

### 3. GET /reproducciones?idUsuario=X
**Función:** Obtener reproducciones de un usuario específico

**Parámetros:**
- `idUsuario`: ID del usuario (query parameter)

**Ejemplo:**
```bash
curl http://localhost:3000/reproducciones?idUsuario=1
```

**Respuesta:** Array de reproducciones del usuario especificado

---

## 🔄 Flujo de Comunicación

### Quién llama a este servidor:

1. **ServidorDeStreaming** (Asíncrono)
   - Cuando un usuario reproduce una canción
   - Envía POST con los metadatos

2. **ServidorDeCalculoPreferencias** (Síncrono)
   - Para calcular preferencias
   - Hace GET por idUsuario

---

## 📁 Estructura de Archivos

```
ServidorDeReproducciones_Correcto/
├── main/
│   └── servidor.go          ← API REST principal
├── go.mod                    ← Sin dependencias externas
├── Reproducciones.json       ← Almacenamiento
└── README.md                 ← Este archivo
```

---

## 💾 Almacenamiento

Los datos se guardan en `Reproducciones.json`:

```json
[
  {
    "id": 1,
    "idUsuario": 1,
    "idCancion": 1,
    "titulo": "Believer",
    "artista": "Imagine Dragons",
    "genero": "Rock",
    "idioma": "Inglés",
    "fechaHora": "2025-11-01T10:30:00Z"
  }
]
```

---

## 🧪 Pruebas

### Almacenar reproducción:
```bash
curl -X POST http://localhost:3000/reproducciones \
  -H "Content-Type: application/json" \
  -d '{"idUsuario":2,"idCancion":3,"titulo":"Test","artista":"Test Artist","genero":"Pop","idioma":"Español"}'
```

### Consultar todas:
```bash
curl http://localhost:3000/reproducciones
```

### Consultar por usuario:
```bash
curl http://localhost:3000/reproducciones?idUsuario=1
```

---

## ✅ Diferencias con la Versión Incorrecta

| Aspecto | Versión Incorrecta | Versión Correcta ✅ |
|---------|-------------------|---------------------|
| Carpeta cliente | ✅ Existe | ❌ NO existe |
| Archivos de audio | ✅ Tiene canciones/ | ❌ NO tiene |
| Streaming | ✅ Implementado | ❌ NO implementado |
| Función | Múltiples | Solo reproducciones |
| Dependencias | Muchas | Ninguna extra |
| Tamaño | 7.6 MB | < 100 KB |

---

## 🎯 Cumplimiento del Requerimiento

✅ **"El servidor de reproducciones permite almacenar una reproducción, y consultar las reproducciones de un id de usuario."**

Esta versión cumple EXACTAMENTE con el requerimiento. Nada más, nada menos.

---

## 🔧 Para Reemplazar en tu Proyecto

1. Elimina la carpeta `ServidorDeReproducciones/` actual
2. Reemplázala con esta carpeta `ServidorDeReproducciones_Correcto/`
3. Renómbrala a `ServidorDeReproducciones/`
4. Ejecuta: `go run main/servidor.go`

---

## 📞 Integración con Otros Componentes

### ServidorDeStreaming debe llamar:
```go
// Después de enviar audio al cliente
http.Post("http://localhost:3000/reproducciones", 
    "application/json", 
    bytes.NewBuffer(jsonData))
```

### ServidorDeCalculoPreferencias debe llamar:
```java
// Usando Feign Client
@RequestLine("GET /reproducciones?idUsuario={idUsuario}")
List<Reproduccion> obtenerReproducciones(@Param("idUsuario") Integer id);
```

---

**Servidor simplificado, funcional y correcto según el requerimiento.** ✅
