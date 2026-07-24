const $ = (id) => document.getElementById(id);
let state = null;
let activeTab = "searchTab";
let theme = localStorage.getItem("multimediasc-theme") || "dark";

function setStatus(text) {
  $("statusPill").textContent = text || "Listo";
  $("statusPillLogin").textContent = text || "Listo";
}

function applyTheme(nextTheme) {
  theme = nextTheme === "light" ? "light" : "dark";
  const root = document.documentElement;
  root.dataset.theme = theme;
  if (theme === "light") {
    root.style.setProperty("--bg", "#edf2f8");
    root.style.setProperty("--panel", "#ffffff");
    root.style.setProperty("--panel-2", "#f6f8fc");
    root.style.setProperty("--line", "rgba(24,39,75,.12)");
    root.style.setProperty("--text", "#122033");
    root.style.setProperty("--muted", "#5b6b7f");
    root.style.setProperty("--accent", "#0f67c2");
    root.style.setProperty("--accent-2", "#1a8f5f");
    root.style.setProperty("--warn", "#9a6b00");
    root.style.setProperty("--shadow", "rgba(18, 32, 51, .12)");
    root.style.setProperty("--input-bg", "#f7f9fd");
    root.style.setProperty("--card-bg", "rgba(18, 32, 51, .03)");
    root.style.setProperty("--card-bg-strong", "rgba(15,103,194,.08)");
  } else {
    root.style.setProperty("--bg", "#0e1621");
    root.style.setProperty("--panel", "#162231");
    root.style.setProperty("--panel-2", "#101a27");
    root.style.setProperty("--line", "rgba(255,255,255,.08)");
    root.style.setProperty("--text", "#e7eef7");
    root.style.setProperty("--muted", "#9fb0c3");
    root.style.setProperty("--accent", "#67b7ff");
    root.style.setProperty("--accent-2", "#7ce6a8");
    root.style.setProperty("--warn", "#ffcc66");
    root.style.setProperty("--shadow", "rgba(0,0,0,.25)");
    root.style.setProperty("--input-bg", "#0f1723");
    root.style.setProperty("--card-bg", "rgba(255,255,255,.03)");
    root.style.setProperty("--card-bg-strong", "rgba(103,183,255,.08)");
  }
  $("themeBtn").textContent = `Tema: ${theme === "light" ? "claro" : "oscuro"}`;
  localStorage.setItem("multimediasc-theme", theme);
}

function planillaClasses(item, active) {
  const cls = ["item"];
  if (active) cls.push("active");
  if (item.error) cls.push("error");
  return cls.join(" ");
}

function badgeClass(item) {
  if (item.error) return "err";
  if (item.selected_docs > 0) return "ok";
  if (item.total_docs > 0) return "warn";
  return "";
}

function showLogin() {
  $("loginScreen").classList.remove("hidden");
  $("mainScreen").classList.add("hidden");
}

function showMain() {
  $("loginScreen").classList.add("hidden");
  $("mainScreen").classList.remove("hidden");
}

function setTab(tabId) {
  activeTab = tabId;
  document.querySelectorAll(".tab").forEach((btn) => {
    btn.classList.toggle("active", btn.dataset.tab === tabId);
  });
  document.querySelectorAll(".tab-panel").forEach((panel) => {
    panel.classList.toggle("hidden", panel.id !== tabId);
  });
}

function canonicalName(name) {
  return String(name || "").trim().replace(/\.[^.]+$/, "");
}

function shortLabel(value) {
  const raw = String(value || "").trim();
  if (!raw) return "";
  const parts = raw.split(/[\\/]/);
  return parts[parts.length - 1] || raw;
}

function usedNamesInPlanilla(planilla) {
  const used = new Set();
  (planilla?.documentos || []).forEach((doc, idx) => {
    const name = canonicalName(doc.nombre);
    if (name && doc.seleccionado) {
      used.add(name.toLowerCase());
    }
  });
  return used;
}

function pickVisibleDoc(planilla, startIndex) {
  const docs = planilla?.documentos || [];
  if (docs.length === 0) return { doc: null, index: -1 };
  for (let i = Math.max(0, startIndex || 0); i < docs.length; i++) {
    if (!docs[i].seleccionado) return { doc: docs[i], index: i };
  }
  for (let i = 0; i < Math.max(0, startIndex || 0); i++) {
    if (!docs[i].seleccionado) return { doc: docs[i], index: i };
  }
  return { doc: null, index: -1 };
}

function renderNameSelect(container, doc, allowedNames, usedNames) {
  const select = document.createElement("select");
  select.className = "name-select";
  const current = canonicalName(doc.nombre);
  const currentKey = current.toLowerCase();
  const available = (allowedNames || []).filter((name) => {
    const key = canonicalName(name).toLowerCase();
    return !usedNames.has(key) || key === currentKey;
  });
  select.innerHTML = `<option value="">Sin nombre</option>`;
  available.forEach((name) => {
    const opt = document.createElement("option");
    opt.value = name;
    opt.textContent = shortLabel(name);
    opt.title = name;
    if (canonicalName(name).toLowerCase() === currentKey) {
      opt.selected = true;
    }
    select.appendChild(opt);
  });
  if (current && !available.some((name) => canonicalName(name).toLowerCase() === currentKey)) {
    const opt = document.createElement("option");
    opt.value = current;
    opt.textContent = `${shortLabel(current)} (activo)`;
    opt.title = current;
    opt.selected = true;
    select.appendChild(opt);
  }
  select.onchange = () => call("RenameDocumento", doc.planillaIndex, doc.index, select.value);
  container.appendChild(select);
}

function render(snapshot) {
  state = snapshot;
  if (!state) return;
  setStatus(state.status);
  if (state.logged_in) {
    showMain();
  } else {
    showLogin();
  }
  $("previewInfo").textContent = state.preview_info || "";
  $("docInfo").textContent = state.doc_info || "";
  $("exportInfo").textContent = state.export_info || "";
  $("batchInfo").textContent = state.batch_info || "";
  $("dirInfo").textContent = state.out_dir ? `Carpeta: ${state.out_dir}` : "Ninguna carpeta seleccionada";
  $("zipInfo").textContent = state.zip_dest ? `ZIP: ${state.zip_dest}` : "El ZIP se generará automáticamente.";
  $("ranges").value = state.range_info ? state.range_info.replace(/^Rangos activos:\s*/, "") : $("ranges").value;
  $("names").value = (state.allowed_names || []).join("\n");

  const planillas = state.planillas || [];
  const list = $("planillasList");
  list.innerHTML = "";
  planillas.forEach((item) => {
    const div = document.createElement("div");
    div.className = planillaClasses(item, item.index === state.active_planilla);
    div.innerHTML = `
      <div class="item-head">
        <strong>${item.error ? `Error ${item.index + 1}` : `${item.tramite}`}</strong>
        <span class="badge ${badgeClass(item)}">${item.error ? 'error' : `${item.selected_docs}/${item.total_docs}`}</span>
      </div>
      <div>${item.error ? item.error : item.paciente || 'Sin paciente'}</div>
      <small>HC: ${item.hc || '-'} · Cédula: ${item.cedula || '-'}${item.error ? '' : ` · Docs: ${item.total_docs}`}</small>
    `;
    div.onclick = () => call('SelectPlanilla', item.index);
    list.appendChild(div);
  });

  const active = planillas[state.active_planilla];
  const docsList = $("docsList");
  docsList.innerHTML = "";
  const basketList = $("basketList");
  basketList.innerHTML = "";
  const basket = [];
  if (active && active.documentos && active.documentos.length > 0) {
    const visible = pickVisibleDoc(active, state.active_documento || 0);
    const doc = visible.doc;
    const usedNames = usedNamesInPlanilla(active);
    if (doc) {
      const div = document.createElement("div");
      div.className = doc.activo ? "doc active single" : "doc single";
      div.innerHTML = `
        <div class="doc-top">
          <div>
            <h3>${doc.descripcion || 'Documento'}</h3>
            <small>${doc.fecha || ''} · ${doc.tipo || ''} · ${doc.estado || ''}</small>
          </div>
          <span class="badge ${doc.seleccionado ? 'ok' : ''}">${doc.seleccionado ? 'Incluido' : 'Pendiente'}</span>
        </div>
        <small>${doc.motivo || ''}</small>
        <div class="doc-meta">Documento ${doc.index + 1} de ${active.documentos.length}</div>
        <div class="name-box"></div>
        <div class="controls">
          <button data-open="open-${doc.planilla}-${doc.index}">Ver</button>
          <button data-toggle="toggle-${doc.planilla}-${doc.index}">${doc.seleccionado ? 'Excluir' : 'Incluir'}</button>
          <button data-mark="mark-${doc.planilla}-${doc.index}">Marcar</button>
        </div>
      `;
      docsList.appendChild(div);

      const nameBox = div.querySelector('.name-box');
      renderNameSelect(nameBox, { ...doc, planillaIndex: active.index }, state.allowed_names || [], usedNames);
      div.querySelector('[data-toggle]').onclick = () => call('ToggleDocumento', active.index, doc.index, !doc.seleccionado);
      div.querySelector('[data-mark]').onclick = () => call('MoveDocumento', doc.index - (state.active_documento || 0));
      div.querySelector('[data-open]').onclick = () => call('OpenDocument', active.index, doc.index);
    } else {
      docsList.innerHTML = '<div class="stats secondary">No hay documentos pendientes para mostrar.</div>';
    }
  } else {
    docsList.innerHTML = '<div class="stats secondary">No hay documentos para mostrar.</div>';
  }

  planillas.forEach((item) => {
    if (!item.documentos) return;
    item.documentos.forEach((doc) => {
      if (!doc.seleccionado) return;
      basket.push({
        planillaIndex: item.index,
        tramite: item.tramite,
        docIndex: doc.index,
        paciente: item.paciente || '',
        descripcion: doc.descripcion || 'Documento',
        nombre: doc.nombre || '',
        tipo: doc.tipo || '',
      });
    });
  });
  if (basket.length === 0) {
    basketList.innerHTML = '<div class="basket-empty">La canasta está vacía. Marca documentos para agregarlos.</div>';
  } else {
    basket.forEach((item) => {
      const div = document.createElement('div');
      div.className = 'basket-item';
      div.innerHTML = `
        <strong>${item.tramite}</strong> · ${item.descripcion}
        <small>${item.nombre || ''}${item.tipo ? ` · ${item.tipo}` : ''}${item.paciente ? ` · ${item.paciente}` : ''}</small>
        <div class="controls">
          <button data-jump>Ir a planilla</button>
          <button data-remove class="ghost">Quitar</button>
        </div>
      `;
      div.querySelector('[data-jump]').onclick = () => {
        call('JumpToPlanilla', item.planillaIndex);
        hideBasket();
      };
      div.querySelector('[data-remove]').onclick = () => {
        call('RemoveSelection', item.planillaIndex, item.docIndex);
      };
      basketList.appendChild(div);
    });
  }
  const basketBtn = $("showBasketBtn");
  if (basketBtn) {
    basketBtn.textContent = `Canasta (${basket.length})`;
  }
  const basketCount = $("basketCount");
  if (basketCount) {
    basketCount.textContent = basket.length === 1 ? "1 documento" : `${basket.length} documentos`;
  }
}

function escapeHtml(value) {
  return String(value)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;');
}

async function call(method, ...args) {
  if (!window.go || !window.go.main || !window.go.main.App || !window.go.main.App[method]) {
    return;
  }
  try {
    const snapshot = await window.go.main.App[method](...args);
    render(snapshot);
  } catch (err) {
    console.error(err);
    setStatus(String(err));
  }
}

function wire() {
  applyTheme(theme);
  $("loginBtn").onclick = () => call('Login', $("user").value, $("pass").value);
  $("logoutBtn").onclick = () => call('Logout');
  $("searchBtn").onclick = () => call('Search');
  $("refreshBtn").onclick = () => call('State');
  $("exportBtn").onclick = () => call('Export');
  $("prevDocBtn").onclick = () => call('MoveDocumento', -1);
  $("nextDocBtn").onclick = () => call('MoveDocumento', 1);
  $("chooseDirBtn").onclick = () => call('SelectOutputDirectory');
  $("clearBasketBtn").onclick = () => call('ClearSelection');
  $("showBasketBtn").onclick = () => showBasket();
  $("themeBtn").onclick = () => applyTheme(theme === "light" ? "dark" : "light");
  $("addAllowedNameBtn").onclick = () => {
    const current = (state?.allowed_names || []);
    const extra = canonicalName($("newAllowedName").value);
    if (!extra) {
      setStatus("ingresa un nombre para agregarlo");
      return;
    }
    const next = [...current];
    if (next.some((name) => canonicalName(name).toLowerCase() === extra.toLowerCase())) {
      setStatus(`el nombre ${extra} ya existe en el catálogo`);
      return;
    }
    next.push(extra);
    $("names").value = next.join("\n");
    $("newAllowedName").value = "";
    call("SetAllowedNames", $("names").value);
    setStatus(`nombre agregado: ${extra}`);
  };
  $("saveNamesBtn").onclick = () => call("SetAllowedNames", $("names").value);
  $("resetNamesBtn").onclick = () => {
    $("names").value = (state?.allowed_names || []).length
      ? "013B\nEpicrisis\nConsentimiento_Informado\nProtocolo_Quirurgico\nResultados_Laboratorio\nImagen_Estudio"
      : $("names").value;
    $("newAllowedName").value = "";
    call("SetAllowedNames", $("names").value);
    setStatus("catálogo restaurado por defecto");
  };
  document.querySelectorAll('[data-close-basket]').forEach((btn) => {
    btn.onclick = () => hideBasket();
  });
  document.querySelectorAll('.tab').forEach((btn) => {
    btn.onclick = () => setTab(btn.dataset.tab);
  });
  $("ranges").onchange = () => call('SetRangeText', $("ranges").value);
  $("names").onchange = () => call('SetAllowedNames', $("names").value);
  $("searchText").onchange = () => call('SetSearchText', $("searchText").value);
  $("searchText").onkeyup = (e) => {
    if (e.key === 'Enter') call('Search');
  };
  $("filter").oninput = () => {
    const q = $("filter").value.toLowerCase().trim();
    document.querySelectorAll('#planillasList .item').forEach((el) => {
      el.style.display = el.textContent.toLowerCase().includes(q) ? '' : 'none';
    });
  };
  setTab(activeTab);
}

function showBasket() {
  $("basketModal").classList.remove("hidden");
  const basketCount = $("basketCount");
  if (basketCount && state) {
    const count = (state.planillas || []).reduce((total, item) => {
      const docs = item.documentos || [];
      return total + docs.filter((doc) => doc.seleccionado).length;
    }, 0);
    basketCount.textContent = count === 1 ? "1 documento" : `${count} documentos`;
  }
}

function hideBasket() {
  $("basketModal").classList.add("hidden");
}

document.addEventListener('DOMContentLoaded', async () => {
  wire();
  await call('State');
});
