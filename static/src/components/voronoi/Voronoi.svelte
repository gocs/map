<script lang="ts">
    import { onMount } from "svelte";
    import { Delaunay, Voronoi } from "d3-delaunay";
    import Cell, { type ICell, type IType } from "./Cell.svelte";

    type point = [number, number, number?, IType?];

    let svg: SVGSVGElement;
    let width: number;
    let height: number;
    let points: point[] = [];
    let voronoi: Voronoi<Delaunay.Point>;
    let delaunay = Delaunay.from(points as [number, number][]);
    let cells = $state<ICell[]>([]);

    async function generatePoints(n: number): Promise<point[]> {
        try {
            const data = await fetch("/land.json");
            const land = await data.json();
            let gen = [] as point[];
            if (land.nodes) {
                for (const node of land.nodes) {
                    gen.push([node.x, node.y, node.id, node.type]);
                }
            }
            return gen;
        } catch (e) {
            console.log("error:", e);
            return Array.from({ length: n }, (_, k) => [
                Math.random() * width,
                Math.random() * height,
                k,
                "sea",
            ]);
        }
    }

    let ws: WebSocket;

    onMount(() => {
        ws = new WebSocket("/ws");
        ws.onopen = function (e) {
            ws.onmessage = function (e) {
                const data = JSON.parse(e.data);
                console.log("message:", data);
            };
            ws.onerror = function (e) {
                console.log("error:", e);
            };
        };

        const resizeObserver = new ResizeObserver(async (entries) => {
            for (let entry of entries) {
                width = entry.contentRect.width;
                height = entry.contentRect.height;
                if (points.length === 0) points = await generatePoints(20);

                drawVoronoi();
            }
        });

        resizeObserver.observe(svg);

        return () => {
            resizeObserver.disconnect();
        };
    });

    function handleDragStart(event: MouseEvent) {
        if (!(event instanceof MouseEvent)) return;
        if (!(event.target instanceof SVGCircleElement)) return;
        const target = event.target;
        if (!target.dataset.index) return;

        const id = parseInt(target.dataset.index);

        function handleDrag(event: MouseEvent) {
            const rect = svg.getBoundingClientRect();
            const x = event.clientX - rect.left;
            const y = event.clientY - rect.top;
            const type = points[id][3];
            points[id] = [x, y, id, type];
            if (ws) {
                const data = { action: "set", id, x, y, type };
                ws.send(JSON.stringify(data));
            }
            drawVoronoi();
        }

        function handleDragEnd() {
            window.removeEventListener("mousemove", handleDrag);
            window.removeEventListener("mouseup", handleDragEnd);
        }

        window.addEventListener("mousemove", handleDrag);
        window.addEventListener("mouseup", handleDragEnd);
    }

    function drawVoronoi() {
        if (!svg || !width || !height) return;

        delaunay = Delaunay.from(points as [number, number][]);
        voronoi = delaunay.voronoi([0, 0, width, height]);

        for (let i = 0; i < points.length; i++) {
            let path = voronoi.cellPolygon(i);
            if (path === null) continue;

            cells[i] = {
                path,
                x: points[i][0],
                y: points[i][1],
                i: points[i][2]!,
                type: points[i][3]!,
            };
        }
    }
</script>

<svg
    bind:this={svg}
    class="w-full h-full"
    onmousedown={handleDragStart}
    role="cell"
    tabindex="-1"
>
    {#each cells as cell, i}
        <Cell bind:cell={cells[i]} debug={false} />
    {/each}
</svg>
