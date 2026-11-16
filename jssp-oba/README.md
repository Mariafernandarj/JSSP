# jssp-oba

Resolución del **Job Shop Scheduling Problem (JSSP)** mediante la heurística **Old Bachelor Acceptance, variante OBA2-N**, escrita en Go. Incluye scripts en Python para graficar el progreso de la búsqueda y el diagrama de Gantt de la mejor solución.

## Requisitos

| Herramienta | Para qué | Versión |
|---|---|---|
| Go | Compilar y ejecutar el solver | La indicada en `go.mod` |
| Python 3 | Gráficas (opcional) | 3.8 o superior |
| `numpy`, `matplotlib` | Gráficas (opcional) | `pip install numpy matplotlib` |

## Estructura

```
jssp-oba/
├── cmd/jssp-oba/main.go        Punto de entrada
├── data/
│   ├── configs/                Archivos de configuración de OBA2-N (JSON)
│   ├── instances/              Instancias de entrada
│   ├── seed/                   Semillas utilizadas
│   └── results/                Resultados generados (txt y csv)
├── internal/
│   ├── problem/                Modelo, parser, validación y cotas LB/UB
│   ├── solution/               Representación del calendario (MachineOrder)
│   ├── objective/              Grafo de precedencias, Kahn, makespan y costo
│   ├── neighborhood/           Generación de vecinos (swap adyacente)
│   ├── oba/                    Algoritmo OBA2-N, configuración y umbral
│   ├── solver/                 Solución inicial y ejecución de experimentos
│   └── reporting/              Exportación de resultados e historial
```

## Compilar

Desde la raíz del proyecto:

```bash
go build -o jssp-oba ./cmd/jssp-oba
```

Esto genera el ejecutable `jssp-oba`. Recompila después de cualquier cambio en el código; un binario viejo no refleja los cambios nuevos.

También puede ejecutarse sin generar el binario:

```bash
go run ./cmd/jssp-oba -instance ... -seeds ... -config ...
```


## Ejecutar

El programa requiere tres argumentos:

| Argumento | Descripción |
|---|---|
| `-instance` | Ruta del archivo de instancia JSSP |
| `-seeds` | Ruta del archivo con las semillas |
| `-config` | Ruta del archivo de configuración JSON |

```bash
~/SeminarioB/JSSP/jssp/JSSP/jssp-oba work*                02:59:22 
❯ go run ./cmd/jssp-oba/main.go \       
    -instance data/instances/100x100.txt \
    -seeds data/seeds/seeds.txt \
    -config data/configs/config100x100.json \
```

Se ejecuta una corrida independiente por cada semilla del archivo. Al iniciar, el programa imprime las cotas inferior (LB) y superior (UB) de la instancia. En consola se muestra por corrida el makespan final, la factibilidad, el porcentaje de aceptación, los puntos de historial y el tiempo.


### Formato de la instancia

```
# Los comentarios empiezan con #
n m
máquina duración  máquina duración  ...   (m pares por cada trabajo, n trabajos)
```

- Primero van `n` (trabajos) y `m` (máquinas).
- Después, `n × m` pares `máquina duración`, trabajo por trabajo y en el orden de sus operaciones.
- Las máquinas se numeran desde **0**.
- Las duraciones deben ser positivas.

### Configuración

Ejemplo de archivo JSON:

```json
{
  "iterations": 1000000,
  "delta": 5,
  "easy_age": 56,
  "history_points": 100
}
```

| Clave | Significado |
|---|---|
| `iterations` | Número de iteraciones por corrida |
| `delta` | Delta del umbral, **en unidades de makespan** (se normaliza internamente por `UB − LB`) |
| `easy_age` | Edad que define una aceptación "fácil" en OBA2-N |
| `history_points` | Cuántos puntos del progreso se guardan por corrida (`0` desactiva el historial) |

Para comparar curvas de varias corridas, todas deben usar el mismo `iterations` y el mismo `history_points`.

## Archivos de salida

Todo se guarda en `data/results/`:

| Archivo | Contenido |
|---|---|
| `resultado_<fecha_hora>.txt` | Solución final, métricas de la mejor solución (makespan, factibilidad, costo, LB, UB, brecha respecto a LB) y datos de la ejecución (iteraciones, aceptados, umbral final, tiempo) |
| `history_<instancia>_seed<N>.csv` | Progreso de la búsqueda de esa semilla |

`<instancia>` es el nombre del archivo de instancia sin extensión; por ejemplo, `JSSP_50_20.txt` produce `history_JSSP_50_20_seed10.csv`.

El historial es un CSV con una línea de metadatos (`# instance=... seed=... LB=... UB=...`) y las columnas:

```
iteration, current_makespan, best_makespan, current_cost, best_cost, threshold, accepted
```

En el `.txt`, la solución aparece como una lista de ternas `[id_operación, máquina, duración]`, recorriendo las máquinas en orden y, dentro de cada una, las operaciones en el orden en que se procesan.

## Gráficas

Los scripts viven en `scripts/` y se ejecutan desde la raíz del proyecto. Si los dejaste en la raíz, quita el prefijo `scripts/` en los comandos.

```bash
pip install numpy matplotlib
```

### Curva de progreso (historial)

Grafica el makespan frente al número de iteraciones, normalizado por LB, con eje X logarítmico. Cada serie se indica como `etiqueta=patrón`; si el patrón coincide con varias semillas se dibuja la media y una banda de mínimo a máximo.

```bash
python3 plot_history.py "OBA2-N=data/results/history_40x40_seed*.csv" -o curva_40x40.png
```

Ejemplo comparando contra una referencia:

```bash
python3 plot_history.py "OBA-N=/home/mafer/jssp-oba/data/results/history_40x40_seed10.csv" -o curva_40x40_10.png
```

Si el script dice `Sin archivos para el patron`, revisa el nombre real con `ls data/results/*.csv`.

### Diagrama de Gantt

Reconstruye el calendario a partir de la instancia y del archivo de resultados, y dibuja una fila por máquina con un color por trabajo.

```bash
python3 scripts/gantt.py \
  data/instances/JSSP_50_20.txt \
  data/results/resultado_20261116_101530.123456789.txt \
  -o gantt_50x20.png
```

Ehmplo:
```
python3 gantt.py /home/mafer/jssp-oba/data/instances/50x20.txt /home/mafer/jssp-oba/data/results/50x20s50.txt -o 50x20_s50.png
```

- `Verificacion OK`: el Gantt corresponde a la solución guardada.
- `AVISO`: no coinciden. Revisa que la instancia sea la misma con la que se hizo la corrida y que la numeración de IDs sea la esperada.
