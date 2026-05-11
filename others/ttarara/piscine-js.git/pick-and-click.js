function pick() {
  const body = document.body;

  const hslDiv = document.createElement('div');
  hslDiv.className = 'hsl';
  const hueDiv = document.createElement('div');
  hueDiv.className = 'hue';
  hueDiv.classList.add('text');
  const lumDiv = document.createElement('div');
  lumDiv.className = 'luminosity';
  lumDiv.classList.add('text');
  body.append(hslDiv, hueDiv, lumDiv);

  const svgNS = 'http://www.w3.org/2000/svg';
  const svg = document.createElementNS(svgNS, 'svg');
  const axisX = document.createElementNS(svgNS, 'line');
  axisX.id = 'axisX';
  const axisY = document.createElementNS(svgNS, 'line');
  body.append(svg);
  axisY.id = 'axisY';
  svg.append(axisX, axisY);

  document.addEventListener('mousemove', e => {
    const x = e.clientX, y = e.clientY;
    const w = window.innerWidth, h = window.innerHeight;
    const hue = Math.round((x / w) * 360);
    const lum = Math.round(100 - (y / h) * 100);
    const hsl = `hsl(${hue}, 50%, ${lum}%)`;

    body.style.background = hsl;

    hslDiv.textContent = hsl;
    hueDiv.textContent = `H: ${hue}`;
    lumDiv.textContent = `L: ${lum}`;

    axisX.setAttribute('x1', x);
    axisX.setAttribute('x2', x);
    axisX.setAttribute('y1', 0);
    axisX.setAttribute('y2', h);
    axisY.setAttribute('y1', y);
    axisY.setAttribute('y2', y);
    axisY.setAttribute('x1', 0);
    axisY.setAttribute('x2', w);
  });

  document.addEventListener('copy', e => {
    const txt = hslDiv.textContent;
    e.clipboardData.setData('text/plain', txt);
    e.preventDefault();
  });

  document.addEventListener('click', () => {
    document.execCommand('copy');
  });
}

export { pick }
