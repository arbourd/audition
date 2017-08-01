import {h, app} from 'hyperapp'

const API_URL = `//${window.location.host}/api`

const Header = () =>
	<div class="audition-header">
		<div class="title">Audition</div>
		<div class="subtitle">A demonstration of Golang, Hyperapp and Terraform</div>
	</div>

const AddMessage = ({state, actions}) =>
	<div class="audition-add-message">
		<div class="subtitle">Add a new message</div>
		<div class="field is-grouped">
			<div class="control is-expanded">
				<input
					class="input"
					type="text"
					placeholder={state.placeholder}
					onkeyup={e => (e.keyCode === 13 ? actions.createMessage() : '')}
					oninput={e => actions.setInput(e.target.value)}
					value={state.input}
				/>
			</div>
			<div class="control">
				<button class="button is-info" onclick={actions.createMessage}>+</button>
			</div>
		</div>
	</div>

const MessageList = ({state, actions}) =>
	<div class="audition-message-list">
		<div class="subtitle">Message List</div>
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
	</div>

const Message = ({id, message, isPalindrome, createdAt, details, actions}) =>
	<div class="audition-message">
		<div class="field is-grouped">
			<div class="control is-expanded">
				<div class="label">{message}</div>
			</div>
			<div class="control">
				<button class="button" onclick={e => actions.setDetailFlag({id: id})}>Details</button>
			</div>
			<div class="control">
				<button class="button is-danger" onclick={e => actions.deleteMessage({id: id})}>-</button>
			</div>
		</div>
		<div class={'audition-details ' + (details === true ? 'show' : 'hide')}>
			<div class="subtitle">Details</div>
			<hr/>
			<p>Palindrome: {isPalindrome.toString()}</p>
			<p>Created: {Date(createdAt).toString()}</p>
		</div>
	</div>

app({
	state: {
		input: '',
		placeholder: '',
		messages: []
	},
	events: {
		init: (state, actions) => (actions.listMessages())
	},
	view: (state, actions) => (
		<div class="container">
			<div class="section">
				<Header />
				<AddMessage state={state} actions={actions}/>
				<MessageList state={state} actions={actions}/>
			</div>
		</div>
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
