//const API_BASE = window.API_BASE || '/routes';
const API_BASE = "/api/routes"
async function loadData() {
    const namespace = document.getElementById('namespace').value.trim();
    const errorEl = document.getElementById('error');
    const loadingEl = document.getElementById('loading');
    const contentEl = document.getElementById('content');
    const refreshBtn = document.getElementById('refreshBtn');

    errorEl.style.display = 'none';
    loadingEl.style.display = 'flex';
    contentEl.style.display = 'none';

    try {
        let url = API_BASE + '?include_deployment=true';
        if (namespace) {
            url += '&namespace=' + encodeURIComponent(namespace);
        }

        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Failed to fetch routes: ' + response.statusText);
        }

        const data = await response.json();
        renderTable(data.routes || []);
        
        refreshBtn.style.display = 'inline-block';
    } catch (err) {
        errorEl.textContent = err.message;
        errorEl.style.display = 'block';
    } finally {
        loadingEl.style.display = 'none';
    }
}

function renderTable(routes) {
    const contentEl = document.getElementById('content');
    
    if (!routes || routes.length === 0) {
        contentEl.innerHTML = '<div class="empty-state">No routes found</div>';
        contentEl.style.display = 'block';
        return;
    }

    let html = `
        <div class="table-wrapper">
            <table>
                <thead>
                    <tr>
                        <th>Route Name</th>
                        <th>Namespace</th>
                        <th>Hostnames</th>
                        <th>Gateways</th>
                        <th>Backend Services</th>
                        <th>Deployments</th>
                    </tr>
                </thead>
                <tbody>
    `;

    routes.forEach(route => {
        const hostnames = route.hostnames?.length 
            ? route.hostnames.map(h => `<span class="tag hostname">${escapeHtml(h)}</span>`).join(' ')
            : '-';
        
        const gateways = route.parentRefs?.length
            ? route.parentRefs.map(g => `<span class="tag gateway">${escapeHtml(g.name)}</span>`).join(' ')
            : '-';

        let servicesHtml = '-';
        let deploymentsHtml = '-';
        
        if (route.rules && route.rules.length > 0) {
            const allBackends = route.rules.flatMap(r => r.backendRefs || []);
            
            if (allBackends.length > 0) {
                servicesHtml = allBackends.map(be => {
                    const port = be.servicePort ? `:${be.servicePort}` : '';
                    return `<span class="tag">${escapeHtml(be.serviceName)}${port} - weight: ${be.weight}</span>`;
                }).join('<br>');

                const allDeployments = allBackends.flatMap(be => be.deploymentInfo || []);
                
                if (allDeployments.length > 0) {
                    deploymentsHtml = `
                        <table class="nested-table">
                            ${allDeployments.map(d => `
                                <tr>
                                    <td><span class="tag">${escapeHtml(d.name)}</span></td>
                                    <td><span class="tag replicas">${d.replicas} replicas</span></td>
                                    <td class="image" title="${escapeHtml(d.image)}">${escapeHtml(d.image)}</td>
                                </tr>
                            `).join('')}
                        </table>
                    `;
                }
            }
        }

        html += `
            <tr>
                <td><strong>${escapeHtml(route.name)}</strong></td>
                <td>${escapeHtml(route.namespace)}</td>
                <td>${hostnames}</td>
                <td>${gateways}</td>
                <td>${servicesHtml}</td>
                <td>${deploymentsHtml}</td>
            </tr>
        `;
    });

    html += '</tbody></table></div>';
    contentEl.innerHTML = html;
    contentEl.style.display = 'block';
}

function escapeHtml(str) {
    if (!str) return '';
    const div = document.createElement('div');
    div.textContent = str;
    return div.innerHTML;
}