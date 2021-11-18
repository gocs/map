let width = window.innerWidth, height = window.innerHeight;
let stroke_width = 5;
let points = 10;

d3.json('land.json')
  .then(graph => {
    let circles = graph.nodes;
    
    let voronoi = d3.Delaunay.from(circles, d => d.x, d => d.y)
      .voronoi([0, 0, width, height]);
    
    // events
    let drag = d3.drag()
      .on("start", (_, d) => circle.filter(p => p === d)
                              .raise()
                              .attr("stroke", "black")
                              .attr("stroke-width", graph.circle.stroke_width))
      .on("drag", (event, d) => (d.x = event.x, d.y = event.y))
      .on("end", (_, d) => {
        fetch('/landjson', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(d),
          }).catch(error => console.error('Error:', error));
        return circle.filter(p => p === d).attr("stroke", "none")
      })
      .on("start.update drag.update end.update", update)
    
    let zoom = d3.zoom()
      .scaleExtent([1, 5])
      .translateExtent([[0, 0], [width, height]])
      .on('zoom', (e, _) => d3.select('svg g').attr('transform', e.transform));
    
    // render
    let map_group = d3.select('svg g');
    
    let mesh = map_group
      .append("path")
      .attr("fill", "black")
      .attr("stroke", graph.line.color)
      .attr("stroke-width", graph.line.stroke_width)
      .attr("d", voronoi.render());
    
    let cell = map_group
      .append("g")
        .attr("fill", "none")
        .attr("pointer-events", "all")
      .selectAll("path")
      .data(circles)
      .join("path")
        .attr("fill", (d, i) => graph.options.colors[d.type])
        .attr("d", (d, i) => voronoi.renderCell(i));
    
    let circle = map_group
      .selectAll("circle")
      .data(circles)
      .join("circle")
        .attr("cx", d => d.x)
        .attr("cy", d => d.y)
        .attr("r", graph.circle.radius)
        .call(drag)
        .attr("fill", "black");
    
    function update() {
      voronoi = d3.Delaunay.from(circles, d => d.x, d => d.y).voronoi([0, 0, width, height]);
      cell
        .attr("fill", (d, _) => graph.options.colors[d.type])
        .attr("d", (_, i) => voronoi.renderCell(i));
      mesh.attr("d", voronoi.render());
      circle.attr("cx", d => d.x).attr("cy", d => d.y);
    }
    
    map_group.node();
    map_group.call(zoom).on("mouseclick");
  });
