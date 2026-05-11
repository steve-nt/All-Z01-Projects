


//Sooo the idea is to fetch the data and make an instance for each hero
//and load them in the prototype Hero
//then render the table with the toRow method

//prototype for each hero with safe property access
function Hero(data) {
    this.name = data.name || 'Unknown';
    this.fullName = data.biography?.fullName || null;
    this.image = data.images?.xs || '';
    this.race = data.appearance?.race || null;
    this.gender = data.appearance?.gender || null;
    this.height = data.appearance?.height || [];
    this.weight = data.appearance?.weight || [];
    this.placeOfBirth = data.biography?.placeOfBirth || null;
    this.alignment = data.biography?.alignment || null;
    this.powerstats = data.powerstats || {};
}

//method to render the hero as a table row
Hero.prototype.toRow = function () {
    // Safe height/weight handling
    const heightDisplay = Array.isArray(this.height) && this.height.length > 0
        ? this.height.join(" / ")
        : "-";
    const weightDisplay = Array.isArray(this.weight) && this.weight.length > 0
        ? this.weight.join(" / ")
        : "-";

    return `
    <tr>
      <td><img class="xs" src="${this.image}" alt="${this.name}"></td>
      <td>${this.name}</td>
      <td>${this.fullName || "-"}</td>
      <td>
        Int: ${this.powerstats.intelligence ?? "-"} |
        Str: ${this.powerstats.strength ?? "-"} |
        Spd: ${this.powerstats.speed ?? "-"} |
        Dur: ${this.powerstats.durability ?? "-"} |
        Pwr: ${this.powerstats.power ?? "-"} |
        Cmb: ${this.powerstats.combat ?? "-"}
      </td>
      <td>${this.race || "-"}</td>
      <td>${this.gender || "-"}</td>
      <td>${this.height.join(" / ")}</td>
      <td>${this.weight.join(" / ")}</td>
      <td>${this.placeOfBirth || "-"}</td>
      <td>${this.alignment || "-"}</td>
    </tr>
  `;
};

// global variable for pagination
let heroObjects = [];
let currentPage = 1;
let pageSize = 20;

let currentSort = { key: "name", order: "asc" };

// helper function so that the totalPages are not NaN as JS cannot convert "all" to number
function getTotalPages() {
    if (pageSize === "all") return 1;
    return Math.ceil(heroObjects.length / pageSize);
}

// render the table with current page
function renderTable() {
    const tableBody = document.getElementById("table-body");
    const startIndex = (currentPage - 1) * pageSize;
    const endIndex = pageSize === "all" ? heroObjects.length : startIndex + pageSize;

    // heroes for current page
    const currentHeroes = heroObjects.slice(startIndex, endIndex);

    tableBody.innerHTML = currentHeroes.map(hero => hero.toRow()).join('');

    updatePaginationInfo();
}

// update pagination controls
function updatePaginationInfo() {
    const pageInfo = document.getElementById("page-info");
    const firstBtn = document.getElementById("first-btn");
    const prevBtn = document.getElementById("prev-btn")
    const nextBtn = document.getElementById("next-btn");
    const lastBtn = document.getElementById("last-btn");


    if (pageSize === "all") {
        pageInfo.textContent = `Showing all ${heroObjects.length} results`;
        firstBtn.disabled = true;
        prevBtn.disabled = true;
        nextBtn.disabled = true;
        lastBtn.disabled = true;
    } else {
        const totalPages = getTotalPages();
        const startItem = (currentPage - 1) * pageSize + 1;
        const endItem = Math.min(currentPage * pageSize, heroObjects.length);

        pageInfo.textContent = `Page ${currentPage} of ${totalPages}`;

        firstBtn.disabled = currentPage === 1;
        prevBtn.disabled = currentPage === 1;
        nextBtn.disabled = currentPage === totalPages;
        lastBtn.disabled = currentPage === totalPages;        
    }
}



// fetch and initialize
fetch("https://rawcdn.githack.com/akabab/superhero-api/0.2.0/api/all.json")
    .then(response => response.json())
    .then(data => {
        heroObjects = data.map(hero => new Hero(hero));//map the raw data to the prototype and store heroes globally

        // intital render
        renderTable();

        document.getElementById("page-size").addEventListener("change", (e) => {
            const value = e.target.value;
            pageSize = value === "all" ? "all" : parseInt(value);
            currentPage = 1; // it resets to the first pagef 
            renderTable();
        });

        // first page
        document.getElementById("first-btn").addEventListener("click", () => {
            if (currentPage > 1) {
                currentPage = 1;
                renderTable();
            }
        });

        // previous button
        document.getElementById("prev-btn").addEventListener("click", () => {
            if (currentPage > 1) {
                currentPage--;
                renderTable();
            }
        });

        // next button
        document.getElementById("next-btn").addEventListener("click", () => {
            const totalPages = getTotalPages();
            if (currentPage < totalPages) {
                currentPage++;
                renderTable();
            }
        });

        // last page
        document.getElementById("last-btn").addEventListener("click", () => {
            const totalPages = getTotalPages();
            if (currentPage < totalPages) {
                currentPage = totalPages;
                renderTable();
            }
        });

    })

    .catch (err => console.error('Error loading data:', err));


// Compact Sorting + Arrows (append-only)
(() => {
    // Map each <th> index to the corresponding Hero property key
    const COLUMNS = [
        null,           // Icon (not sortable)
        "name",
        "fullName",
        "powerstats",   // sortable by sum of stats
        "race",
        "gender",
        "height",
        "weight",
        "placeOfBirth",
        "alignment",
    ];

    // Extractors for each key so we avoid switch/case duplication
    const EXTRACT = {
        name: h => h.name,
        fullName: h => h.fullName,
        race: h => h.race,
        gender: h => h.gender,
        placeOfBirth: h => h.placeOfBirth,
        alignment: h => h.alignment,
        height: h => h.height,     // array → normalized
        weight: h => h.weight,     // array → normalized
        powerstats: h => {
            const ps = h.powerstats;
            if (!ps) return null;
            let sum = 0, seen = 0;
            for (const k of ["intelligence", "strength", "speed", "durability", "power", "combat"]) {
                const n = typeof ps[k] === "number" ? ps[k] : parseFloat(ps[k]);
                if (Number.isFinite(n)) { sum += n; seen++; }
            }
            return seen ? sum : null;
        },
    };

    // Normalizer:
    // - missing → null
    // - arrays → metric value if available ("cm"/"kg"), otherwise first element
    // - "78 kg" → 78
    // - fallback → lowercase string
    const normalize = v => {
        if (v == null || v === "" || v === "-" || (Array.isArray(v) && v.length === 0)) return null;
        if (Array.isArray(v)) {
            const metric = v.find(x => x && /cm|kg/i.test(String(x)));
            v = metric || v[0];
            if (!v) return null;
        }
        const s = String(v);
        const hasDigit = /\d/.test(s);
        const n = parseFloat(s.replace(/[^0-9.+-]/g, ""));
        if (hasDigit && Number.isFinite(n)) return n;
        return s.toLowerCase();
    };

    // Comparator factory
    const makeComparator = (key, order) => (a, b) => {
        const rawA = (EXTRACT[key] ?? (() => null))(a);
        const rawB = (EXTRACT[key] ?? (() => null))(b);
        const va = normalize(rawA), vb = normalize(rawB);

        // Missing values always last (regardless of asc/desc)
        const am = va === null, bm = vb === null;
        if (am && bm) return 0;
        if (am) return 1;
        if (bm) return -1;

        if (va < vb) return order === "asc" ? -1 : 1;
        if (va > vb) return order === "asc" ? 1 : -1;
        return 0;
    };

    // Attach header click listeners only once
    function ensureHeaderUI() {
        const ths = document.querySelectorAll("thead th");
        ths.forEach((th, i) => {
            const key = COLUMNS[i];
            if (!key) return; // skip Icon
            th.classList.add("sortable");
            th.dataset.sortKey = key;
            th.style.cursor = "pointer";
            th.title = "Click to sort";
            th.onclick = () => {
                if (currentSort.key === key) {
                    currentSort.order = currentSort.order === "asc" ? "desc" : "asc";
                } else {
                    currentSort.key = key;
                    currentSort.order = "asc";
                }
                heroObjects.sort(makeComparator(currentSort.key, currentSort.order));
                currentPage = 1;
                renderTable();
                updateHeaderArrows();
            };
        });
    }

    // Update arrow indicators on headers
    function updateHeaderArrows() {
        const ths = document.querySelectorAll("thead th");
        ths.forEach(th => th.classList.remove("sorted-asc", "sorted-desc"));
        const active = document.querySelector(`thead th[data-sort-key="${currentSort.key}"]`);
        if (active) active.classList.add(currentSort.order === "asc" ? "sorted-asc" : "sorted-desc");
    }

    // Initialization without touching your existing fetch/render code:
    // 1) Attach header UI after DOM is ready
    // 2) Wait until heroObjects is populated, then apply default sort by Name asc
    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", ensureHeaderUI, { once: true });
    } else {
        ensureHeaderUI();
    }

    const wait = setInterval(() => {
        if (Array.isArray(heroObjects) && heroObjects.length) {
            clearInterval(wait);
            // Default sort by name asc if not already done
            if (currentSort?.key !== "name" || currentSort?.order !== "asc") {
                currentSort.key = "name";
                currentSort.order = "asc";
                heroObjects.sort(makeComparator("name", "asc"));
                currentPage = 1;
                renderTable();
            }
            updateHeaderArrows();
        }
    }, 50);
})();

// Minimal Interactive Search (append-only) 
// Filters by name on each keystroke. Keeps current sort. Missing other code unchanged.
(() => {
    let MASTER = null; // snapshot of full dataset

    // Filter function (case-insensitive substring match on name)
    function applySearch(q) {
        if (!MASTER && Array.isArray(heroObjects)) MASTER = heroObjects.slice();
        const needle = String(q || "").toLowerCase().trim();

        heroObjects = needle
            ? MASTER.filter(h => String(h.name || "").toLowerCase().includes(needle))
            : MASTER.slice();

        // Re-apply current sorting if available; fallback to name ASC
        if (typeof makeComparator === "function" && currentSort?.key) {
            heroObjects.sort(makeComparator(currentSort.key, currentSort.order || "asc"));
        } else {
            heroObjects.sort((a, b) => String(a.name).localeCompare(String(b.name)));
        }

        currentPage = 1;
        renderTable();
        if (typeof updateHeaderArrows === "function") updateHeaderArrows();
    }

    // Bind once the input event
    function bind() {
        const el = document.getElementById("search");
        if (!el) return;
        el.addEventListener("input", e => applySearch(e.target.value));
    }

    // Attach after DOM is ready
    if (document.readyState === "loading") {
        document.addEventListener("DOMContentLoaded", bind, { once: true });
    } else {
        bind();
    }

    // Create MASTER once data is loaded (no changes to your fetch)
    const t = setInterval(() => {
        if (Array.isArray(heroObjects) && heroObjects.length) {
            clearInterval(t);
            MASTER = heroObjects.slice();
        }
    }, 50);
})();
