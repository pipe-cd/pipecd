/*
 * TOC Scroll Spy
 * Highlights the current section in the table of contents based on scroll position
 */

(function () {
  'use strict';

  // Wait for DOM to be ready
  document.addEventListener('DOMContentLoaded', function () {
    const toc = document.querySelector('.td-toc');
    if (!toc) {
      return; // No TOC on this page
    }

    // Get all TOC links
    const tocLinks = toc.querySelectorAll('a[href^="#"]');
    if (tocLinks.length === 0) {
      return; // No anchor links in TOC
    }

    // Get all section headings
    const sections = [];
    tocLinks.forEach(function (link) {
      const id = link.getAttribute('href').substring(1);
      const section = document.getElementById(id);
      if (section) {
        sections.push({
          id: id,
          element: section,
          link: link,
        });
      }
    });

    if (sections.length === 0) {
      return; // No matching sections found
    }

    // Function to get the current active section
    function getActiveSection() {
      const scrollPosition = window.scrollY + 100; // Offset for better UX

      // Find the section that's currently in view
      for (let i = sections.length - 1; i >= 0; i--) {
        const section = sections[i];
        if (section.element.offsetTop <= scrollPosition) {
          return section;
        }
      }

      // If we're at the top of the page, return the first section
      return sections[0];
    }

    // Function to update active state
    function updateActiveLink() {
      const activeSection = getActiveSection();

      // Remove active class from all links
      sections.forEach(function (section) {
        section.link.classList.remove('active');
        section.link.style.fontWeight = '';
        section.link.style.color = '';
      });

      // Add active class to current link
      if (activeSection) {
        activeSection.link.classList.add('active');
        activeSection.link.style.fontWeight = 'bold';
        activeSection.link.style.color = '#007bff'; // Bootstrap primary color
      }
    }

    // Throttle function to limit scroll event frequency
    let scrollTimeout;
    function throttledUpdate() {
      if (scrollTimeout) {
        window.cancelAnimationFrame(scrollTimeout);
      }
      scrollTimeout = window.requestAnimationFrame(function () {
        updateActiveLink();
      });
    }

    // Listen for scroll events
    window.addEventListener('scroll', throttledUpdate);

    // Initial update
    updateActiveLink();

    // Also update when hash changes (clicking on TOC links)
    window.addEventListener('hashchange', function () {
      setTimeout(updateActiveLink, 100);
    });
  });
})();
