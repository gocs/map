<script lang="ts" module>
    export const color = {
        cliff: "#55555560",
        forest: "#11ff1160",
        shore: "#f6e3ab60",
        sea: "#31c1ff60",
    };
    export type IType = "cliff" | "forest" | "shore" | "sea";

    export interface ICell {
        path: Delaunay.Polygon;
        x: number;
        y: number;
        i: number;
        type: IType;
    }
</script>

<script lang="ts">
    import type { Delaunay } from "d3-delaunay";
    let { cell = $bindable(), debug = false }: { cell: ICell; debug: boolean } =
        $props();

</script>

{#if cell}
    <path d="M{cell.path.join('L')}Z" fill={color[cell.type]} stroke="#fff" />
    <g transform="translate({cell.x} {cell.y})">
        <circle r="5" fill="#fff" cursor="move" data-index={cell.i} />
        {#if debug}<text
                x="10"
                y="10"
                class="font-mono text-xs z-10"
                fill="#fff">{Math.round(cell.x)},{Math.round(cell.y)}</text
            >{/if}
    </g>
{/if}
