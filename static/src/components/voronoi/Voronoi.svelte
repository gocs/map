<script lang="ts">
    import { onMount } from "svelte";
    import { Delaunay } from "d3-delaunay";
    let svg: SVGSVGElement;
    let width: number;
    let height: number;
    let points: [number, number][] = [];
    let voronoi: any;

    type terrain = "cliff" | "forrest" | "shore" | "sea";

    const color = {
        cliff: "#55555560",
        forest: "#11ff1160",
        shore: "#f6e3ab60",
        sea: "#31c1ff60",
    };

    async function generatePoints(n: number): [number, number][] {
        const data = await fetch("/land.json");
        const land = await data.json();
        let gen = [];
        if (land.nodes) {
            for (const node of land.nodes) {
                gen.push([node.x, node.y, node.id, node.type]);
            }
        }
        return gen;
    }

    function mksvg(shape: string) {
        return document.createElementNS("http://www.w3.org/2000/svg", shape);
    }
    let ws;

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
                if (points.length === 0) {
                    points = await generatePoints(200);
                }
                drawVoronoi();
            }
        });

        resizeObserver.observe(svg);

        return () => {
            resizeObserver.disconnect();
        };
    });

    function handleDragStart(event: MouseEvent) {
        const target = event.target as SVGCircleElement;
        if (!target.dataset.index) return;

        const id = parseInt(target.dataset.index);

        function handleDrag(event: MouseEvent) {
            const rect = svg.getBoundingClientRect();
            const x = event.clientX - rect.left;
            const y = event.clientY - rect.top;
            console.log("points[id]:", points[id]);
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

        svg.innerHTML = "";
        const delaunay = Delaunay.from(points);
        voronoi = delaunay.voronoi([0, 0, width, height]);

        for (let i = 0; i < points.length; i++) {
            const cell = voronoi.cellPolygon(i);
            if (!cell) return;
            const [cx, cy, id, type] = points[i];

            const g = mksvg("g");
            const path = mksvg("path");
            path.setAttribute("d", `M${cell.join("L")}Z`);
            path.setAttribute("fill", color[type]);
            path.setAttribute("stroke", "#00000000");
            g.appendChild(path);

            const circle = mksvg("circle");
            circle.setAttribute("cx", `${cx}`);
            circle.setAttribute("cy", `${cy}`);
            circle.setAttribute("r", "5");
            circle.setAttribute("fill", "#ffffff99");
            circle.setAttribute("cursor", "move");
            circle.dataset.index = `${i}`;
            g.dataset.index = `${i}`;
            g.appendChild(circle);

            // const text = mksvg("text");
            // text.setAttribute("x", `${cx + 10}`);
            // text.setAttribute("y", `${cy - 10}`);
            // text.setAttribute("font-family", "monospace");
            // text.setAttribute("font-size", "12");
            // text.setAttribute("fill", "#333");
            // text.setAttribute("class", "z-index:1");
            // text.textContent = `${Math.round(cx)},${Math.round(cy)}`;
            // g.appendChild(text);
            svg.appendChild(g);
        }
    }
</script>

<svg bind:this={svg} class="w-full h-full" on:mousedown={handleDragStart}></svg>
