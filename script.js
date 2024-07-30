document.addEventListener('DOMContentLoaded', function() {
  // Dropdown functionality
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

  // Fetch and populate data on page load
  const fetchData = async () => {
    try {
      // Fetch learner data
      const learnerResponse = await fetch('http://localhost:3000/learners/1');
      const learnerData = await learnerResponse.json();
      document.getElementById('learnerName').textContent = learnerData.name;
      document.getElementById('learnerId').textContent = learnerData.id;
      document.getElementById('completedSessions').textContent = learnerData.completedSessions;
      document.getElementById('upcomingSessions').textContent = learnerData.upcomingSessions;
      document.getElementById('totalHours').textContent = learnerData.totalHours;
      document.getElementById('averageRating').textContent = learnerData.averageRating;
      document.getElementById('walletBalance').textContent = `KES ${learnerData.walletBalance}`;

      // Fetch and populate courses data
      const coursesResponse = await fetch('http://localhost:3000/courses');
      const courses = await coursesResponse.json();
      const coursesTableBody = document.getElementById('coursesTableBody');
      coursesTableBody.innerHTML = '';
      courses.forEach(course => {
        const row = coursesTableBody.insertRow();
        row.insertCell(0).textContent = course.name;
        row.insertCell(1).textContent = course.tutor;
        const progressCell = row.insertCell(2);
        progressCell.innerHTML = `
          <div class="progress-bar">
            <div class="progress" style="width: ${course.progress}%">${course.progress}%</div>
          </div>
        `;
        row.insertCell(3).textContent = new Date(course.nextSession).toLocaleString();
        const actionCell = row.insertCell(4);
        const viewButton = document.createElement('button');
        viewButton.textContent = 'View Details';
        viewButton.className = 'action-btn view-btn';
        viewButton.onclick = () => viewCourseDetails(course);
        actionCell.appendChild(viewButton);
      });

      // Fetch and populate tutors data
      const tutorsResponse = await fetch('http://localhost:3000/tutors');
      const tutors = await tutorsResponse.json();
      const tutorsTableBody = document.getElementById('tutorsTableBody');
      tutorsTableBody.innerHTML = '';
      tutors.forEach(tutor => {
        const row = tutorsTableBody.insertRow();
        row.insertCell(0).textContent = tutor.name;
        row.insertCell(1).textContent = tutor.subject;
        row.insertCell(2).textContent = tutor.rating;
        row.insertCell(3).textContent = tutor.rate;
        const actionCell = row.insertCell(4);
        const bookButton = document.createElement('button');
        bookButton.textContent = 'Book Session';
        bookButton.className = 'action-btn book-btn';
        bookButton.onclick = () => bookSession(tutor);
        actionCell.appendChild(bookButton);
      });

    } catch (error) {
      console.error('Error fetching data:', error);
    }
  };

  // Call fetchData to populate the data
  fetchData();

  // Handle form submission
  const form = document.getElementById('registration-form');
  form.addEventListener('submit', async (e) => {
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
    try {
      const response = await fetch('/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(user)
      });
      const data = await response.json();
      console.log(data);
    } catch (error) {
      console.error('Error registering user:', error);
    }
  });

  // Additional functions
  function viewCourseDetails(course) {
    alert(`Viewing details for ${course.name}`);
    // Implement course details view
  }

  function bookSession(tutor) {
    alert(`Booking a session with ${tutor.name}`);
    // Implement booking functionality
  }
});
