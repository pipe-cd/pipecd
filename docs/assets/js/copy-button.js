/**
 * Copy button functionality for code snippets
 * Adds a copy button to all highlighted code blocks
 */

(function() {
  'use strict';

  // SVG icon for copy button
  const copyIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path><rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect></svg>`;
  
  const checkIcon = `<svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>`;

  /**
   * Add copy button to a code block
   */
  function addCopyButton(codeBlock) {
    // Skip if button already added
    if (codeBlock.parentElement.querySelector('.copy-button')) {
      return;
    }

    // Get the text content from the code block
    const codeText = codeBlock.textContent;

    // Create button container
    const buttonContainer = document.createElement('div');
    buttonContainer.className = 'copy-button-container';

    // Create button
    const button = document.createElement('button');
    button.className = 'copy-button';
    button.title = 'Copy code';
    button.innerHTML = copyIcon;
    button.type = 'button';

    // Add click handler
    button.addEventListener('click', function(e) {
      e.preventDefault();
      
      // Copy to clipboard
      if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(codeText).then(() => {
          // Show success state
          showCopySuccess(button);
        }).catch(() => {
          // Fallback: try using execCommand
          fallbackCopy(codeText, button);
        });
      } else {
        // Fallback for older browsers
        fallbackCopy(codeText, button);
      }
    });

    buttonContainer.appendChild(button);
    codeBlock.parentElement.insertBefore(buttonContainer, codeBlock);
  }

  /**
   * Fallback copy method for older browsers
   */
  function fallbackCopy(text, button) {
    const textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.style.position = 'fixed';
    textarea.style.left = '-9999px';
    document.body.appendChild(textarea);
    
    try {
      textarea.select();
      document.execCommand('copy');
      showCopySuccess(button);
    } catch (err) {
      console.error('Failed to copy text:', err);
    } finally {
      document.body.removeChild(textarea);
    }
  }

  /**
   * Show copy success feedback
   */
  function showCopySuccess(button) {
    const originalHTML = button.innerHTML;
    button.classList.add('copied');
    button.innerHTML = checkIcon;
    button.title = 'Copied!';

    // Reset after 2 seconds
    setTimeout(() => {
      button.classList.remove('copied');
      button.innerHTML = originalHTML;
      button.title = 'Copy code';
    }, 2000);
  }

  /**
   * Initialize copy buttons when DOM is ready
   */
  function initCopyButtons() {
    // Target code blocks within .highlight divs (Hugo/Chroma output)
    const codeBlocks = document.querySelectorAll('.highlight > pre > code');
    
    codeBlocks.forEach(codeBlock => {
      addCopyButton(codeBlock);
    });
  }

  // Initialize when DOM is ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initCopyButtons);
  } else {
    initCopyButtons();
  }

  // Also support dynamic content (if needed in future)
  if (window.MutationObserver) {
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.addedNodes.length) {
          mutation.addedNodes.forEach((node) => {
            if (node.nodeType === 1) { // Element node
              const codeBlocks = node.querySelectorAll?.('.highlight > pre > code') || [];
              codeBlocks.forEach(addCopyButton);
            }
          });
        }
      });
    });

    observer.observe(document.body, {
      childList: true,
      subtree: true
    });
  }
})();
