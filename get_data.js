document.addEventListener('DOMContentLoaded', function() {
    // Fetch learner data
    fetch('http://localhost:3000/learners/1')
      .then(response => response.json())
      .then(data => {
        document.getElementById('learnerName').textContent = data.name;
        document.getElementById('learnerId').textContent = data.id;
        document.getElementById('completedSessions').textContent = data.completedSessions;
        document.getElementById('upcomingSessions').textContent = data.upcomingSessions;
        document.getElementById('totalHours').textContent = data.totalHours;
        document.getElementById('averageRating').textContent = data.averageRating;
        document.getElementById('walletBalance').textContent = `KES ${data.walletBalance}`;
      });
  
    // Fetch courses data
    fetch('http://localhost:3000/courses')
      .then(response => response.json())
      .then(courses => {
        const tableBody = document.getElementById('coursesTableBody');
        tableBody.innerHTML = '';
        courses.forEach(course => {
          const row = tableBody.insertRow();
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
      });
  
    // Fetch tutors data
    fetch('http://localhost:3000/tutors')
      .then(response => response.json())
      .then(tutors => {
        const tableBody = document.getElementById('tutorsTableBody');
        tableBody.innerHTML = '';
        tutors.forEach(tutor => {
          const row = tableBody.insertRow();
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
      });
  });
  