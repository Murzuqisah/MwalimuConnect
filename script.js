document.addEventListener('DOMContentLoaded', function() {
    const dropdownLinks = document.querySelectorAll('.dropdown-content a');
    const hiddenInput = document.getElementById('roles');
    
    dropdownLinks.forEach(link => {
        link.addEventListener('click', function(event) {
            event.preventDefault(); // Prevents the default anchor behavior
            const value = this.getAttribute('data-value');
            hiddenInput.value = value; // Set the hidden input value
            document.querySelector('.dropbtn').textContent = this.textContent; // Update button text
        });
    });
});

// script.js

const form = document.getElementById('registration-form');

form.addEventListener('submit', (e) => {
  e.preventDefault();

  const username = document.getElementById('username').value;
  const email = document.getElementById('email').value;
  const password = document.getElementById('password').value;

  // Create a new user object
  const user = {
    username,
    email,
    password
  };

  // Send a POST request to the server to store the user details
  fetch('/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(user)
  })
  .then(response => response.json())
  .then(data => console.log(data))
  .catch(error => console.error(error));
});