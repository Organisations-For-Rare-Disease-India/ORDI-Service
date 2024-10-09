document.addEventListener('DOMContentLoaded', function () {
    // Handle dropdown button clicks
    document.querySelectorAll('.dropdown-button').forEach(button => {
      button.addEventListener('click', function () {
        const menu = this.closest('.dropdown').querySelector('.dropdown-menu');
        // Toggle the 'hidden' class to show/hide the menu
        menu.classList.toggle('hidden');
      });
    });
  
    // Close the dropdown if clicked outside
    window.addEventListener('click', function (event) {
      if (!event.target.matches('.dropdown-button') && !event.target.closest('.dropdown-menu')) {
        document.querySelectorAll('.dropdown-menu').forEach(menu => {
          menu.classList.add('hidden');
        });
      }
    });
  });
  