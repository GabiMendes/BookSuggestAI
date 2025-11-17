// Dashboard JavaScript functionality
let selectedSpreadsheetId = '';

document.addEventListener('DOMContentLoaded', function() {
    loadSpreadsheets();
});

async function loadSpreadsheets() {
    try {
        const response = await fetch('/api/v1/spreadsheets');
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Erro ao carregar planilhas');
        }
        
        const select = document.getElementById('spreadsheet-select');
        select.innerHTML = '<option value="">Selecione uma planilha...</option>';
        
        // Verificar se spreadsheets existe e é um array
        if (data.spreadsheets && Array.isArray(data.spreadsheets)) {
            data.spreadsheets.forEach(sheet => {
                const option = document.createElement('option');
                option.value = sheet.id;
                option.textContent = sheet.name;
                select.appendChild(option);
            });
        } else {
            showAlert('Nenhuma planilha encontrada no Google Drive.', 'warning');
        }
    } catch (error) {
        console.error('Erro ao carregar planilhas:', error);
        showAlert('Erro ao carregar suas planilhas. Verifique se você tem planilhas no Google Drive.', 'danger');
    }
}

async function analyzeSpreadsheet() {
    const select = document.getElementById('spreadsheet-select');
    selectedSpreadsheetId = select.value;
    
    if (!selectedSpreadsheetId) {
        showAlert('Por favor, selecione uma planilha primeiro.', 'warning');
        return;
    }
    
    showLoading(true);
    
    try {
        const response = await fetch(`/api/v1/spreadsheets/${selectedSpreadsheetId}/analyze`);
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Erro na análise');
        }
        
        displayAnalysis(data.analysis);
        showRecommendationButton();
    } catch (error) {
        console.error('Erro na análise:', error);
        showAlert(`Erro ao analisar a planilha: ${error.message}`, 'danger');
    } finally {
        showLoading(false);
    }
}

function displayAnalysis(analysis) {
    const content = document.getElementById('analysis-content');
    
    // Garantir que os arrays existem para evitar erros de null
    const favoriteGenres = analysis.favorite_genres || [];
    const favoriteAuthors = analysis.favorite_authors || [];
    
    content.innerHTML = `
        <div class="row">
            <div class="col-md-3">
                <div class="analysis-stat">
                    <i class="fas fa-book text-primary"></i>
                    <h4>${analysis.total_books}</h4>
                    <p class="mb-0">Livros Lidos</p>
                </div>
            </div>
            <div class="col-md-3">
                <div class="analysis-stat">
                    <i class="fas fa-star text-warning"></i>
                    <h4>${analysis.average_rating.toFixed(1)}</h4>
                    <p class="mb-0">Nota Média</p>
                </div>
            </div>
            <div class="col-md-3">
                <div class="analysis-stat">
                    <i class="fas fa-tags text-info"></i>
                    <h4>${favoriteGenres.length}</h4>
                    <p class="mb-0">Gêneros Favoritos</p>
                </div>
            </div>
            <div class="col-md-3">
                <div class="analysis-stat">
                    <i class="fas fa-user-edit text-success"></i>
                    <h4>${favoriteAuthors.length}</h4>
                    <p class="mb-0">Autores Preferidos</p>
                </div>
            </div>
        </div>
        
        <div class="row mt-3">
            <div class="col-md-6">
                <h6><i class="fas fa-heart text-danger me-2"></i>Gêneros Favoritos:</h6>
                <p class="text-muted">${favoriteGenres.length > 0 ? favoriteGenres.join(', ') : 'Nenhum padrão identificado'}</p>
            </div>
            <div class="col-md-6">
                <h6><i class="fas fa-pen text-primary me-2"></i>Autores Preferidos:</h6>
                <p class="text-muted">${favoriteAuthors.length > 0 ? favoriteAuthors.join(', ') : 'Nenhum padrão identificado'}</p>
            </div>
        </div>
        
        <div class="alert alert-info mt-3">
            <i class="fas fa-lightbulb me-2"></i>
            <strong>Resumo:</strong> ${analysis.reading_patterns}
        </div>
    `;
    
    document.getElementById('analysis-section').classList.remove('d-none');
    document.getElementById('analysis-section').classList.add('fade-in');
}

function showRecommendationButton() {
    document.getElementById('recommendation-button').classList.remove('d-none');
    document.getElementById('recommendation-button').classList.add('fade-in');
}

async function generateRecommendations() {
    if (!selectedSpreadsheetId) {
        showAlert('Nenhuma planilha selecionada.', 'warning');
        return;
    }
    
    showLoading(true);
    
    try {
        const response = await fetch('/api/v1/recommendations', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                spreadsheet_id: selectedSpreadsheetId
            })
        });
        
        const data = await response.json();
        
        if (!response.ok) {
            throw new Error(data.error || 'Erro ao gerar recomendações');
        }
        
        displayRecommendations(data.recommendations);
    } catch (error) {
        console.error('Erro ao gerar recomendações:', error);
        showAlert(`Erro ao gerar recomendações: ${error.message}`, 'danger');
    } finally {
        showLoading(false);
    }
}

function displayRecommendations(recommendations) {
    const container = document.getElementById('recommendations-cards');
    
    // Verificar se recommendations existe e é um array
    if (!recommendations || !Array.isArray(recommendations)) {
        showAlert('Erro: Nenhuma recomendação recebida', 'warning');
        return;
    }
    
    container.innerHTML = recommendations.map((rec, index) => `
        <div class="col-md-4 mb-4">
            <div class="card recommendation-card">
                <div class="card-header bg-primary text-white">
                    <h5 class="mb-0">
                        <i class="fas fa-trophy me-2"></i>Recomendação ${rec.position}
                    </h5>
                </div>
                <div class="card-body">
                    <h5 class="card-title text-primary">${rec.title}</h5>
                    <h6 class="card-subtitle mb-2 text-muted">
                        <i class="fas fa-user me-1"></i>${rec.author}
                    </h6>
                    <p class="card-text">${rec.justification}</p>
                    
                    <div class="mt-3">
                        <h6>Esta recomendação foi útil?</h6>
                        <div class="btn-group" role="group">
                            <button type="button" class="btn btn-outline-success btn-sm" onclick="raterecommendation(${rec.id}, 5)">
                                <i class="fas fa-thumbs-up"></i> Sim
                            </button>
                            <button type="button" class="btn btn-outline-danger btn-sm" onclick="rateRecommendation(${rec.id}, 1)">
                                <i class="fas fa-thumbs-down"></i> Não
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    `).join('');
    
    document.getElementById('recommendations-section').classList.remove('d-none');
    document.getElementById('recommendations-section').classList.add('fade-in');
    
    // Scroll to recommendations
    document.getElementById('recommendations-section').scrollIntoView({ 
        behavior: 'smooth' 
    });
}

function rateRecommendation(recommendationId, rating) {
    // TODO: Implement recommendation rating
    console.log(`Rating recommendation ${recommendationId} with ${rating} stars`);
    showAlert('Obrigado pelo seu feedback!', 'success');
}

function showLoading(show) {
    if (show) {
        document.getElementById('loading-section').classList.remove('d-none');
    } else {
        document.getElementById('loading-section').classList.add('d-none');
    }
}

function showAlert(message, type = 'info') {
    // Remove existing alerts
    const existingAlerts = document.querySelectorAll('.alert');
    existingAlerts.forEach(alert => alert.remove());
    
    const alertDiv = document.createElement('div');
    alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
    alertDiv.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;
    
    // Insert at the top of the container
    const container = document.querySelector('.container-fluid');
    container.insertBefore(alertDiv, container.firstChild);
    
    // Auto-remove after 5 seconds
    setTimeout(() => {
        if (alertDiv.parentNode) {
            alertDiv.remove();
        }
    }, 5000);
}

async function logout() {
    try {
        const response = await fetch('/auth/logout', {
            method: 'POST'
        });
        
        if (response.ok) {
            window.location.href = '/';
        } else {
            showAlert('Erro ao fazer logout', 'danger');
        }
    } catch (error) {
        console.error('Erro no logout:', error);
        showAlert('Erro ao fazer logout', 'danger');
    }
}