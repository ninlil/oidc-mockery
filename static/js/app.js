// OIDC Mockery JavaScript

document.addEventListener('DOMContentLoaded', function() {
    console.log('OIDC Mockery JavaScript loaded');
    
    // Add click animations to persona cards
    const personaCards = document.querySelectorAll('.persona-card');
    
    personaCards.forEach(card => {
        card.addEventListener('click', function() {
            // Find the button within this card and click it
            const button = this.querySelector('button[type="submit"]');
            if (button) {
                button.click();
            }
        });
        
        // Add hover effect
        card.addEventListener('mouseenter', function() {
            this.style.transform = 'translateY(-2px)';
            this.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.15)';
        });
        
        card.addEventListener('mouseleave', function() {
            this.style.transform = 'translateY(0)';
            this.style.boxShadow = 'none';
        });
    });
    
    // Handle form submission with loading state
    const forms = document.querySelectorAll('form');
    const buttons = document.querySelectorAll('button[type="submit"]');
    
    console.log('Found', forms.length, 'forms and', buttons.length, 'submit buttons');
    
    // Store original button text
    buttons.forEach(button => {
        button.setAttribute('data-original-text', button.textContent);
    });
    
    forms.forEach(form => {
        console.log('Setting up form submission handler for form:', form.action, form.method);
        form.addEventListener('submit', function(event) {
            console.log('Form submission detected!', this.action, this.method);
            
            // Find the submit button that was clicked
            const submitButton = event.submitter || this.querySelector('button[type="submit"]');
            
            if (submitButton) {
                console.log('Disabling button:', submitButton.textContent);
                
                // Set the persona_id from the button's data attribute
                const personaId = submitButton.getAttribute('data-persona-id');
                if (personaId) {
                    const personaInput = document.getElementById('persona-id-input');
                    if (personaInput) {
                        personaInput.value = personaId;
                        console.log('Set persona_id to:', personaId);
                    }
                }
                
                submitButton.disabled = true;
                submitButton.textContent = 'Loading...';
            }
            
            // Allow the form to submit normally
            console.log('Allowing form to submit...');
        });
    });
    
    // Console info for developers
    console.log('OIDC Mockery UI loaded');
    console.log('Available personas:', document.querySelectorAll('.persona-card').length);
});
