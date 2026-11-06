#!/usr/bin/env python3
"""
Genera un diagrama de Gantt a partir del archivo de resultados de jssp-oba.

Uso:
    python3 gantt.py instancia.txt resultado.txt [-o gantt.png]

Requiere: matplotlib  (pip install matplotlib)

Supuestos (se verifican al final contra el "Makespan" del archivo):
  * Cada elemento de la lista de solucion es [x, y, z] y x es el ID de la
    operacion. Los otros dos valores se ignoran: los tiempos se recalculan.
  * ID de operacion = trabajo * m + indice_en_el_trabajo (numeracion por filas
    de la instancia, 0-based).
  * Las operaciones de cada maquina aparecen en el archivo en el orden en que
    se procesan (MachineOrder).
"""
import argparse
import colorsys
import re
import sys
from collections import defaultdict

import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from matplotlib.patches import Patch


def read_instance(path):
    """Formato: n m, luego n*m pares (maquina duracion). '#' inicia comentario."""
    with open(path, encoding="utf-8") as f:
        text = re.sub(r"#.*", "", f.read())
    nums = [int(t) for t in re.findall(r"-?\d+", text)]
    n, m = nums[0], nums[1]
    pairs = nums[2:]
    if len(pairs) != 2 * n * m:
        sys.exit(f"Instancia inconsistente: se esperaban {2*n*m} enteros y hay {len(pairs)}.")
    machine, duration = {}, {}
    for j in range(n):
        for k in range(m):
            op_id = j * m + k
            machine[op_id] = pairs[2 * op_id]
            duration[op_id] = pairs[2 * op_id + 1]
    return n, m, machine, duration


def read_solution(path):
    with open(path, encoding="utf-8") as f:
        text = f.read()
    start = text.find("Solucion")
    end = text.find("Informacion", start if start >= 0 else 0)
    block = text[start:end] if start >= 0 else text
    triples = re.findall(r"\[\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\]", block)
    op_ids = [int(a) for a, _, _ in triples]
    mk = re.search(r"Makespan:\s*(\d+)", text)
    return op_ids, int(mk.group(1)) if mk else None


def build_schedule(n, m, machine, duration, op_ids):
    """Reconstruye inicios y fines (calendario semiactivo) con Kahn."""
    order = defaultdict(list)
    for op in op_ids:
        order[machine[op]].append(op)
    if sorted(op_ids) != list(range(n * m)):
        sys.exit("La solucion no contiene cada operacion exactamente una vez.")

    preds = defaultdict(list)
    succs = defaultdict(list)
    for j in range(n):
        for k in range(1, m):
            a, b = j * m + k - 1, j * m + k
            succs[a].append(b); preds[b].append(a)
    for ops in order.values():
        for a, b in zip(ops, ops[1:]):
            succs[a].append(b); preds[b].append(a)

    indeg = {op: len(preds[op]) for op in range(n * m)}
    start = {op: 0 for op in range(n * m)}
    queue = [op for op, d in indeg.items() if d == 0]
    done = 0
    while queue:
        u = queue.pop()
        done += 1
        end_u = start[u] + duration[u]
        for v in succs[u]:
            start[v] = max(start[v], end_u)
            indeg[v] -= 1
            if indeg[v] == 0:
                queue.append(v)
    if done != n * m:
        sys.exit("El calendario tiene ciclos: no se puede dibujar un Gantt.")
    return start, order


def job_colors(n):
    cols = []
    for j in range(n):
        h = (j * 0.618033988749895) % 1.0          # tonos bien separados
        cols.append(colorsys.hsv_to_rgb(h, 0.55, 0.92))
    return cols


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("instancia")
    ap.add_argument("resultado")
    ap.add_argument("-o", "--output", default="gantt.png")
    ap.add_argument("--dpi", type=int, default=150)
    args = ap.parse_args()

    n, m, machine, duration = read_instance(args.instancia)
    op_ids, reported = read_solution(args.resultado)
    start, order = build_schedule(n, m, machine, duration, op_ids)
    makespan = max(start[o] + duration[o] for o in range(n * m))

    if reported is not None and reported != makespan:
        print(f"AVISO: makespan recalculado = {makespan}, pero el archivo dice {reported}. "
              "Revisa la instancia usada o el significado de los campos.", file=sys.stderr)
    else:
        print(f"Verificacion OK: makespan = {makespan}")

    colors = job_colors(n)
    n_mach = max(machine.values()) + 1
    fig_h = max(4, 0.32 * n_mach + 2)
    fig, ax = plt.subplots(figsize=(18, fig_h))

    for mach, ops in order.items():
        for op in ops:
            j = op // m
            ax.barh(mach, duration[op], left=start[op], height=0.8,
                    color=colors[j], edgecolor="black", linewidth=0.3)
            if duration[op] >= makespan / 90:        # etiqueta solo si cabe
                ax.text(start[op] + duration[op] / 2, mach, str(j),
                        ha="center", va="center", fontsize=5)

    ax.set_xlim(0, makespan)
    ax.set_ylim(-0.6, n_mach - 0.4)
    ax.invert_yaxis()
    ax.set_yticks(range(n_mach))
    ax.set_yticklabels([f"M{i}" for i in range(n_mach)], fontsize=7)
    ax.set_xlabel("Tiempo")
    ax.set_title(f"Diagrama de Gantt - {n} trabajos x {m} operaciones - makespan = {makespan}")
    ax.grid(axis="x", linestyle=":", alpha=0.5)
    if n <= 20:
        ax.legend(handles=[Patch(color=colors[j], label=f"Trabajo {j}") for j in range(n)],
                  loc="upper left", bbox_to_anchor=(1.005, 1), fontsize=7, ncol=1)
    fig.tight_layout()
    fig.savefig(args.output, dpi=args.dpi)
    print(f"Guardado: {args.output}")


if __name__ == "__main__":
    main()
