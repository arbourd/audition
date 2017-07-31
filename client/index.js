import {h, app} from 'hyperapp'

const API_URL = `//${window.location.host}/api`

// Message component
const Message = ({id, message, isPalindrome, createdAt, details, actions}) =>
	<li>
		{message}
		<button onclick={e => actions.setDetailFlag({id: id})}>Details</button>
		<button onclick={e => actions.deleteMessage({id: id})}>-</button>
		<div class={details === true ? 'show' : 'hide'}>
			Palindrome: {isPalindrome.toString()} <br />
			Created: {Date(createdAt).toString()}
		</div>
	</li>

app({
	state: {
		input: '',
		placeholder: 'Add new message',
		messages: []
	},
	events: {
		init: (state, actions) => (actions.listMessages())
	},
	view: (state, actions) => (
		<main>
			<h1>Messages</h1>
			<input
				type="text"
				placeholder={state.placeholder}
				onkeyup={e => (e.keyCode === 13 ? actions.createMessage() : '')}
				oninput={e => actions.setInput(e.target.value)}
				value={state.input}
				autofocus
			/>
			<button onclick={actions.createMessage}>+</button>
			<ul>
				{state.messages.map(m =>
					<Message
						id={m.id}
						message={m.message}
						isPalindrome={m.isPalindrome}
						createdAt={m.createdAt}
						details={m.details}
						actions={actions}
					/>
				)}
			</ul>
		</main>
	),
	actions: {
		// Helper actions for manipulating state
		setInput: (state, actions, input) => ({input}),
		setMessages: (state, actions, messages) => ({messages}),
		setDetailFlag: (state, actions, {id}) => {
			for (const msg of state.messages) {
				if (msg.id === id) {
					msg.details = !msg.details
				}
			}
			return {messages: state.messages}
		},

		// HTTP services
		listMessages: (state, actions) => {
			fetch(`${API_URL}/messages`)
			.then(res => res.json())
			.then(messages => {
				if (messages.error) {
					return handleError(messages)
				}

				actions.setInput('')
				actions.setMessages(processMessages(messages))
			})
			.catch(err => console.log('Error: ' + err.message))
		},
		createMessage: (state, actions) => {
			fetch(`${API_URL}/messages`, {
				method: 'POST',
				body: JSON.stringify({message: state.input})
			})
			.then(res => res.json())
			.then(message => {
				if (message.error) {
					return handleError(message)
				}

				actions.setInput('')
				actions.setMessages(processMessages(state.messages.concat(message)))
			})
			.catch(err => console.log('Error: ' + err.message))
		},
		deleteMessage: (state, actions, {id}) => {
			fetch(`${API_URL}/messages/${id}`, {method: 'DELETE'})
			.then(res => {
				if (res.status !== 204) {
					return res.json()
				}
				actions.setMessages(state.messages.filter(m => (m.id !== id)))
			})
			.then(error => {
				if (error) {
					handleError(error)
				}
			})
			.catch(err => console.log('Error: ' + err.message))
		}
	}
})

function handleError(err) {
	alert(`${err.error}: ${err.message}`)
}

// Adds a default `details: false` to message's state
function processMessages(msgs) {
	if (!Array.isArray(msgs)) {
		msgs = [msgs]
	}

	for (const msg of msgs) {
		if (!('details' in msg)) {
			msg.details = false
		}
	}
	return msgs
}
