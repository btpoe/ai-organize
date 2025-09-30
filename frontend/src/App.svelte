<script>
  import { AnalyzeDirectory, ExecuteMoves } from '../wailsjs/go/main/App';

  let directoryPath = '';
  let analyzing = false;
  let executing = false;
  let analysisResult = null;
  let proposedMoves = [];
  let selectedMoves = new Set();
  let executeResult = null;
  let errorMessage = '';

  // Group moves by category for better UI
  $: groupedMoves = proposedMoves.reduce((acc, move, idx) => {
    if (!acc[move.category]) {
      acc[move.category] = [];
    }
    acc[move.category].push({ ...move, index: idx });
    return acc;
  }, {});

  async function analyzeDirectory() {
    if (!directoryPath.trim()) {
      errorMessage = 'Please enter a directory path';
      return;
    }

    analyzing = true;
    errorMessage = '';
    executeResult = null;

    try {
      const result = await AnalyzeDirectory(directoryPath);

      if (result.error) {
        errorMessage = result.error;
        analysisResult = null;
        proposedMoves = [];
      } else {
        analysisResult = result;
        proposedMoves = result.proposedMoves;
        // Select all moves by default
        selectedMoves = new Set(proposedMoves.map((_, idx) => idx));
      }
    } catch (error) {
      errorMessage = 'Error analyzing directory: ' + error.message;
    } finally {
      analyzing = false;
    }
  }

  function toggleMove(index) {
    if (selectedMoves.has(index)) {
      selectedMoves.delete(index);
    } else {
      selectedMoves.add(index);
    }
    selectedMoves = selectedMoves; // Trigger reactivity
  }

  function toggleCategory(category) {
    const categoryMoves = groupedMoves[category];
    const allSelected = categoryMoves.every(m => selectedMoves.has(m.index));

    if (allSelected) {
      // Deselect all in category
      categoryMoves.forEach(m => selectedMoves.delete(m.index));
    } else {
      // Select all in category
      categoryMoves.forEach(m => selectedMoves.add(m.index));
    }
    selectedMoves = selectedMoves;
  }

  async function executeMoves() {
    const movesToExecute = proposedMoves.filter((_, idx) => selectedMoves.has(idx));

    if (movesToExecute.length === 0) {
      errorMessage = 'No moves selected';
      return;
    }

    executing = true;
    errorMessage = '';

    try {
      const result = await ExecuteMoves(movesToExecute);
      executeResult = result;

      // Clear the analysis after successful execution
      if (result.success > 0) {
        // Remove executed moves from the list
        proposedMoves = proposedMoves.filter((_, idx) => !selectedMoves.has(idx));
        selectedMoves.clear();

        if (proposedMoves.length === 0) {
          analysisResult = null;
        }
      }
    } catch (error) {
      errorMessage = 'Error executing moves: ' + error.message;
    } finally {
      executing = false;
    }
  }

  function editDestination(index) {
    const newPath = prompt('Enter new destination path:', proposedMoves[index].destinationPath);
    if (newPath && newPath !== proposedMoves[index].destinationPath) {
      proposedMoves[index].destinationPath = newPath;
      proposedMoves = proposedMoves; // Trigger reactivity
    }
  }

  function reset() {
    directoryPath = '';
    analysisResult = null;
    proposedMoves = [];
    selectedMoves = new Set();
    executeResult = null;
    errorMessage = '';
  }
</script>

<main>
  <div class="container">
    <header>
      <h1>📁 File Organizer</h1>
      <p>Automatically organize your files into logical folders</p>
    </header>

    <div class="input-section">
      <div class="input-group">
        <input
          type="text"
          bind:value={directoryPath}
          placeholder="Enter directory path (e.g., /Users/username/Downloads)"
          disabled={analyzing || executing}
        />
        <button
          class="btn-primary"
          on:click={analyzeDirectory}
          disabled={analyzing || executing}
        >
          {analyzing ? 'Analyzing...' : 'Analyze Directory'}
        </button>
      </div>
      {#if errorMessage}
        <div class="error-message">{errorMessage}</div>
      {/if}
    </div>

    {#if analysisResult && proposedMoves.length > 0}
      <div class="results-section">
        <div class="results-header">
          <h2>
            Analysis Results: {analysisResult.totalFiles} files found,
            {proposedMoves.length} moves proposed
          </h2>
          <div class="results-actions">
            <span>{selectedMoves.size} selected</span>
            <button class="btn-secondary" on:click={reset}>Reset</button>
            <button
              class="btn-success"
              on:click={executeMoves}
              disabled={executing || selectedMoves.size === 0}
            >
              {executing ? 'Executing...' : `Execute ${selectedMoves.size} Moves`}
            </button>
          </div>
        </div>

        <div class="categories">
          {#each Object.entries(groupedMoves) as [category, moves]}
            <div class="category-section">
              <div class="category-header">
                <label class="checkbox-label">
                  <input
                    type="checkbox"
                    checked={moves.every(m => selectedMoves.has(m.index))}
                    on:change={() => toggleCategory(category)}
                  />
                  <h3>{category} ({moves.length})</h3>
                </label>
              </div>

              <div class="moves-list">
                {#each moves as move}
                  <div class="move-item" class:selected={selectedMoves.has(move.index)}>
                    <label class="checkbox-label">
                      <input
                        type="checkbox"
                        checked={selectedMoves.has(move.index)}
                        on:change={() => toggleMove(move.index)}
                      />
                      <div class="move-details">
                        <div class="file-name">{move.fileName}</div>
                        <div class="file-paths">
                          <div class="path">
                            <span class="label">From:</span>
                            <span class="value">{move.sourcePath}</span>
                          </div>
                          <div class="path">
                            <span class="label">To:</span>
                            <span class="value">{move.destinationPath}</span>
                          </div>
                          <div class="reason">{move.reason}</div>
                        </div>
                      </div>
                    </label>
                    <button
                      class="btn-edit"
                      on:click={() => editDestination(move.index)}
                      title="Edit destination"
                    >
                      ✏️
                    </button>
                  </div>
                {/each}
              </div>
            </div>
          {/each}
        </div>
      </div>
    {:else if analysisResult && proposedMoves.length === 0}
      <div class="info-message">
        ✅ All files are already organized! No moves needed.
      </div>
    {/if}

    {#if executeResult}
      <div class="execute-results">
        <h3>Execution Results</h3>
        <div class="stats">
          <div class="stat success">
            <span class="label">Successfully moved:</span>
            <span class="value">{executeResult.success}</span>
          </div>
          {#if executeResult.failed > 0}
            <div class="stat failure">
              <span class="label">Failed:</span>
              <span class="value">{executeResult.failed}</span>
            </div>
          {/if}
        </div>

        {#if executeResult.createdFolders.length > 0}
          <div class="created-folders">
            <strong>Created folders:</strong>
            <ul>
              {#each executeResult.createdFolders as folder}
                <li>{folder}</li>
              {/each}
            </ul>
          </div>
        {/if}

        {#if executeResult.failedFiles.length > 0}
          <div class="failed-files">
            <strong>Failed operations:</strong>
            <ul>
              {#each executeResult.failedFiles as error}
                <li>{error}</li>
              {/each}
            </ul>
          </div>
        {/if}
      </div>
    {/if}
  </div>
</main>

<style>
  :global(body) {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    min-height: 100vh;
  }

  main {
    padding: 2rem;
    min-height: 100vh;
  }

  .container {
    max-width: 1200px;
    margin: 0 auto;
    background: white;
    border-radius: 16px;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
    padding: 2rem;
  }

  header {
    text-align: center;
    margin-bottom: 2rem;
    padding-bottom: 1.5rem;
    border-bottom: 2px solid #f0f0f0;
  }

  h1 {
    color: #333;
    font-size: 2.5rem;
    margin-bottom: 0.5rem;
  }

  header p {
    color: #666;
    font-size: 1.1rem;
  }

  .input-section {
    margin-bottom: 2rem;
  }

  .input-group {
    display: flex;
    gap: 1rem;
    margin-bottom: 1rem;
  }

  input[type="text"] {
    flex: 1;
    padding: 0.875rem 1rem;
    border: 2px solid #e0e0e0;
    border-radius: 8px;
    font-size: 1rem;
    transition: border-color 0.2s;
  }

  input[type="text"]:focus {
    outline: none;
    border-color: #667eea;
  }

  input[type="text"]:disabled {
    background: #f5f5f5;
    cursor: not-allowed;
  }

  button {
    padding: 0.875rem 1.5rem;
    border: none;
    border-radius: 8px;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
  }

  button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .btn-primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  .btn-primary:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
  }

  .btn-secondary {
    background: #f0f0f0;
    color: #333;
  }

  .btn-secondary:hover:not(:disabled) {
    background: #e0e0e0;
  }

  .btn-success {
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
    color: white;
  }

  .btn-success:hover:not(:disabled) {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(17, 153, 142, 0.4);
  }

  .btn-edit {
    padding: 0.5rem;
    background: transparent;
    font-size: 1.2rem;
    min-width: auto;
  }

  .btn-edit:hover {
    background: #f0f0f0;
  }

  .error-message {
    padding: 1rem;
    background: #fee;
    color: #c33;
    border-radius: 8px;
    border-left: 4px solid #c33;
  }

  .info-message {
    padding: 2rem;
    text-align: center;
    background: #e8f5e9;
    color: #2e7d32;
    border-radius: 8px;
    font-size: 1.2rem;
  }

  .results-section {
    margin-top: 2rem;
  }

  .results-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 2px solid #f0f0f0;
  }

  .results-header h2 {
    color: #333;
    font-size: 1.5rem;
  }

  .results-actions {
    display: flex;
    gap: 1rem;
    align-items: center;
  }

  .results-actions span {
    color: #666;
    font-weight: 600;
  }

  .categories {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .category-section {
    border: 2px solid #f0f0f0;
    border-radius: 8px;
    overflow: hidden;
  }

  .category-header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 1rem 1.5rem;
  }

  .category-header h3 {
    color: white;
    font-size: 1.2rem;
    margin: 0;
  }

  .category-header .checkbox-label {
    color: white;
  }

  .moves-list {
    padding: 1rem;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

  .move-item {
    display: flex;
    gap: 1rem;
    padding: 1rem;
    background: #fafafa;
    border-radius: 6px;
    border: 2px solid transparent;
    transition: all 0.2s;
  }

  .move-item.selected {
    background: #f0f4ff;
    border-color: #667eea;
  }

  .checkbox-label {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    cursor: pointer;
    flex: 1;
  }

  .checkbox-label input[type="checkbox"] {
    margin-top: 0.25rem;
    width: 18px;
    height: 18px;
    cursor: pointer;
  }

  .move-details {
    flex: 1;
  }

  .file-name {
    font-weight: 600;
    color: #333;
    margin-bottom: 0.5rem;
    font-size: 1.05rem;
  }

  .file-paths {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .path {
    display: flex;
    gap: 0.5rem;
    font-size: 0.9rem;
  }

  .path .label {
    color: #666;
    font-weight: 600;
    min-width: 50px;
  }

  .path .value {
    color: #333;
    word-break: break-all;
  }

  .reason {
    color: #999;
    font-size: 0.85rem;
    margin-top: 0.25rem;
    font-style: italic;
  }

  .execute-results {
    margin-top: 2rem;
    padding: 1.5rem;
    background: #f9f9f9;
    border-radius: 8px;
    border: 2px solid #e0e0e0;
  }

  .execute-results h3 {
    margin-bottom: 1rem;
    color: #333;
  }

  .stats {
    display: flex;
    gap: 2rem;
    margin-bottom: 1rem;
  }

  .stat {
    display: flex;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    border-radius: 6px;
    font-weight: 600;
  }

  .stat.success {
    background: #e8f5e9;
    color: #2e7d32;
  }

  .stat.failure {
    background: #ffebee;
    color: #c62828;
  }

  .created-folders,
  .failed-files {
    margin-top: 1rem;
  }

  .created-folders ul,
  .failed-files ul {
    margin-top: 0.5rem;
    padding-left: 1.5rem;
  }

  .created-folders li {
    color: #666;
    padding: 0.25rem 0;
  }

  .failed-files {
    color: #c62828;
  }

  .failed-files li {
    padding: 0.25rem 0;
  }
</style>