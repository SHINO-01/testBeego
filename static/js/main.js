document.addEventListener('DOMContentLoaded', function () {
    // State management
    let currentImageId = null;

    // DOM Elements
    const tabs = document.querySelectorAll('.tab-btn');
    const tabContents = document.querySelectorAll('.tab-pane');
    const loadingOverlay = document.querySelector('.loading-overlay');

    // Toast notification function
    function showToast(message, type = 'success') {
        const toast = document.createElement('div');
        toast.className = `toast ${type}`;
        toast.textContent = message;
        document.body.appendChild(toast);

        setTimeout(() => {
            toast.classList.add('show');
        }, 100);

        setTimeout(() => {
            toast.classList.remove('show');
            setTimeout(() => {
                document.body.removeChild(toast);
            }, 300);
        }, 3000);
    }

    // Loading state management
    function setLoading(isLoading) {
        if (isLoading) {
            loadingOverlay.classList.add('active');
        } else {
            loadingOverlay.classList.remove('active');
        }
    }

    // Tab switching
    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            const targetTab = tab.dataset.tab;

            tabs.forEach(t => t.classList.remove('active'));
            tabContents.forEach(c => c.classList.remove('active'));

            tab.classList.add('active');
            document.getElementById(targetTab).classList.add('active');

            if (targetTab === 'voting') loadNewCatToVote();
            if (targetTab === 'breeds') loadBreedsDropdown();
            if (targetTab === 'favs') loadFavorites();
            if (targetTab === 'vote_history') loadVoteHistory();
        });
    });

    // Voting functionality
    async function loadNewCatToVote() {
        try {
            setLoading(true);
            const response = await fetch('/api/cats/random');
            const [data] = await response.json();

            currentImageId = data.id;

            const catImage = document.querySelector('.cat-image');
            catImage.innerHTML = `<img src="${data.url}" alt="Cat">`;
        } catch (error) {
            showToast('Error loading cat image', 'error');
            console.error('Error:', error);
        } finally {
            setLoading(false);
        }
    }

    //fav button
    document.querySelector('.favorite-btn').addEventListener('click', async () => {
        const payload = {
            image_id: currentImageId // Replace with the actual image ID
        };

        try {
            const response = await fetch('/api/favorites', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json', // Ensure JSON content type
                },
                body: JSON.stringify(payload), // Serialize payload to JSON
            });

            if (response.ok) {
                showToast('Added to favorites!');
                loadNewCatToVote();
            } else {
                const errorData = await response.json();
                console.error('Error:', errorData);
                showToast('Failed to add favorite', 'error');
            }
        } catch (error) {
            console.error('Error adding favorite:', error);
            showToast('Error adding favorite', 'error');
        }
    });

    // Vote buttons
    document.querySelectorAll('.vote-btn').forEach(btn => {
        btn.addEventListener('click', async () => {
            if (!currentImageId) {
                console.error('No image ID available for voting');
                return;
            }

            const payload = {
                image_id: currentImageId,
                value: parseInt(btn.dataset.vote) // Ensure this is an integer
            };

            console.log('Payload:', payload);

            try {
                const response = await fetch('/api/vote', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(payload) // Serialize the payload
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    console.error('Error response:', errorData);
                    showToast('Failed to vote', 'error');
                } else {
                    showToast('Vote Recorded Successfully!');
                    loadNewCatToVote();
                    const responseData = await response.json();
                    console.log('Success response:', responseData);
                }
            } catch (error) {
                console.error('Fetch error:', error);
            }
        });
    });

    // Breed dropdown and details functionality
    async function loadBreedsDropdown() {
        try {
            const response = await fetch('/api/breeds');
            const breeds = await response.json();
            const dropdown = document.getElementById('breed-select');

            dropdown.innerHTML = '<option value="" disabled selected>Select a Breed</option>'; // Reset dropdown
            breeds.forEach(breed => {
                const option = document.createElement('option');
                option.value = breed.id;
                option.textContent = breed.name;
                dropdown.appendChild(option);
            });

            dropdown.addEventListener('change', async (e) => {
                const breedId = e.target.value;
                const breed = breeds.find(b => b.id === breedId);
                if (breed) displayBreedDetails(breed);
            });
        } catch (error) {
            showToast('Error loading breeds', 'error');
            console.error('Error:', error);
        }
    }

    async function displayBreedDetails(breed) {
        document.getElementById('breed-name').textContent = breed.name;
        document.getElementById('breed-description').textContent = breed.description;
        document.getElementById('breed-origin').textContent = breed.origin;
        document.getElementById('breed-temperament').textContent = breed.temperament;
        const wikiLink = document.getElementById('breed-wiki');
        if (breed.wikipedia_url) {
            wikiLink.href = breed.wikipedia_url;
            wikiLink.style.display = 'inline';
        } else {
            wikiLink.style.display = 'none';
        }

        const slideshow = document.getElementById('slideshow-images');
        slideshow.innerHTML = 'Loading images...';

        try {
            const params = new URLSearchParams({
                breed_ids: breed.id,  // e.g. "beng"
                limit: 8
            });

            const response = await fetch(`/api/breed-images?${params.toString()}`);
            if (!response.ok) throw new Error('Failed to fetch images from server');

            const images = await response.json();
            if (!images.length) throw new Error('No images found for this breed');

            slideshow.innerHTML = images
                .map(img => `<img src="${img.url}" alt="${breed.name}" class="slideshow-img">`)
                .join('');
            startSlideshow();

            console.log(`Loaded ${images.length} images for breed: ${breed.name}`);
        } catch (error) {
            slideshow.innerHTML = `<p>Error loading images: ${error.message}</p>`;
            console.error(error);
        }

        document.getElementById('breed-details').style.display = 'block';
    }

    let slideshowInterval;
    function startSlideshow() {
        const images = document.querySelectorAll('.slideshow-img');
        const dotContainer = document.createElement('div');
        dotContainer.className = 'slideshow-dots';
        const slideshowContainer = document.querySelector('.slideshow-container');

        // Remove existing dots
        const existingDots = document.querySelector('.slideshow-dots');
        if (existingDots) existingDots.remove();

        // Create dot indicators
        images.forEach((_, index) => {
            const dot = document.createElement('span');
            dot.className = 'slideshow-dot';
            if (index === 0) dot.classList.add('active');
            dot.addEventListener('click', () => {
                setActiveSlide(index);
            });
            dotContainer.appendChild(dot);
        });

        slideshowContainer.appendChild(dotContainer);

        let currentIndex = 0;

        function setActiveSlide(index, direction = 'next') {
            const currentSlide = images[currentIndex];
            const nextSlide = images[index];

            // Remove any ongoing transitions
            images.forEach(img => {
                if (img !== currentSlide && img !== nextSlide) {
                    img.style.transition = 'none';
                    img.classList.remove('active');
                }
            });

            // Enable transitions for current and next slides
            currentSlide.style.transition = '';
            nextSlide.style.transition = '';

            // Update active states
            currentSlide.classList.remove('active');
            nextSlide.classList.add('active');

            // Update dots
            const dots = document.querySelectorAll('.slideshow-dot');
            dots.forEach((dot, i) => {
                dot.classList.toggle('active', i === index);
            });

            currentIndex = index;
        }

        // Automatically transition slides every 5 seconds
        if (slideshowInterval) clearInterval(slideshowInterval);
        slideshowInterval = setInterval(() => {
            const nextIndex = (currentIndex + 1) % images.length;
            setActiveSlide(nextIndex);
        }, 3000);
    }

    // Favorites functionality
    async function loadFavorites() {
        try {
            const response = await fetch('/api/favorites');
            const favorites = await response.json();

            const container = document.querySelector('.favorites-grid');
            if (favorites.length === 0) {
                container.innerHTML = '<p class="no-favorites">No favorites yet. Start adding some cats!</p>';
                return;
            }

            container.innerHTML = favorites.map(fav => `
                <div class="favorite-card" data-id="${fav.id}">
                    <img src="${fav.image.url}" alt="Favorite cat">
                    <div class="favorite-details">
                        <span class="favorite-date">${new Date(fav.created_at).toLocaleString('en-US', {
                            dateStyle: 'short',
                            timeStyle: 'short'
                          })}</span>
                        <button class="remove-btn" onclick="removeFavorite(${fav.id})">&#10005 Remove</button>
                    </div>             
                </div>
            `).join('');
        } catch (error) {
            showToast('Error loading favorites', 'error');
            console.error('Error:', error);
        }
    }

    window.removeFavorite = async function (id) {
        try {
            const response = await fetch(`/api/favorites/${id}`, {
                method: 'DELETE'
            });

            if (response.ok) {
                showToast('Removed from favorites');
                loadFavorites();
            } else {
                throw new Error('Failed to remove favorite');
            }
        } catch (error) {
            showToast('Failed to remove from favorites', 'error');
            console.error('Error:', error);
        }
    };

    async function loadVoteHistory() {
        try {
            setLoading(true); // Show loading overlay if desired
    
            const response = await fetch('/api/vote_history');
            if (!response.ok) {
                throw new Error('Failed to fetch vote history');
            }
    
            // The response should be an array of votes
            const votes = await response.json();
    
            const container = document.querySelector('.history_grid');
            if (!votes || votes.length === 0) {
                container.innerHTML = '<p class="no-history">No voting history available.</p>';
                return;
            }

            // Build grid cards
            container.innerHTML = votes.map(vote => {
                // 1 = upvote => green border, -1 = downvote => red border
                const borderColor = vote.value === 1 ? 'green' : 'red';
    
                // The Cat API's votes often include an `image` field with `url`
                // If not present, use a placeholder or handle gracefully
                const imageUrl = vote.image && vote.image.url
                    ? vote.image.url
                    : '/static/img/placeholder.png';
    
                // Date formatting
                const dateTime = new Date(vote.created_at).toLocaleString('en-US', {
                    dateStyle: 'short',
                    timeStyle: 'short'
                  });
    
                return `
                    <div class="vote-history-card" style="border: 6px solid ${borderColor};">
                        <img src="${imageUrl}" alt="Voted cat">
                        <div class="vote-history-details">
                            <!-- First row: Date and Vote ID -->
                            <div class="detail-item">
                                <span class="vote-history-date">${dateTime}</span>
                            </div>
                            <div class="detail-item">
                                <p>Vote ID: ${vote.id}</p>
                            </div>
                        </div>
                    </div>
                `;
            }).join('');
    
        } catch (error) {
            showToast('Error loading vote history', 'error');
            console.error('Error:', error);
        } finally {
            setLoading(false); // Hide loading overlay
        }
    }
      
    // Initialize the app
    loadNewCatToVote();
});
