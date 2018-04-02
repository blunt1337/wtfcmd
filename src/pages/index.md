<banner></banner>

---

<h1 menu-ignore>What's the fu*king command!</h1>

You want to run a command and you don't remember it?
You have co-workers asking you for commands every time?
You have a ton of shell scripts and don't remember what each one is for?

Build a single json/yaml file with all the commands available in your project, and anyone working with you, know in seconds how to run everything.

## Features
- Run commands with 'wtf my_command_name [params...]', all defined in a json/yaml file.
- Auto-generated --help! so you know what they are for and how to use them.
- Autocomplete all of your commands and flags!
- Possibility to have custom commands for bash and powershell.exe, of course.
- So much more, like parameter checks, built in functions to print errors, etc.

---

## Getting started

### Install the command
The following instructions are for **{{ install_mode }}**. Install on <template v-for="mode in other_install_modes"><a href="#install-the-command" @click.prevent="install_mode = mode">{{ mode }}? </a></template>

<install-go v-if="install_mode === 'anywhere with Go compiler'"></install-go><install-windows v-if="install_mode === 'Windows'"></install-windows><install-mac v-if="install_mode === 'Mac'"></install-mac><install-linux v-if="install_mode === 'Linux'"></install-linux>

### Setup your project

It's so simple, just create a .wtfcmd.yml file in the root of your project.
You can define all the commands the project needs in there, check the [samples](/samples) page to get inspired.

<script>
import Banner from 'js/components/banner'
import InstallWindows from './installs/windows'
import InstallMac from './installs/mac'
import InstallLinux from './installs/linux'
import InstallGo from './installs/go'

export default {
	data: () => ({
		install_mode: null,
	}),
	mounted() {
		let app = navigator.appVersion
		if (app.indexOf('Win') != -1) {
			this.install_mode = 'Windows'
		} else if (app.indexOf('Mac') != -1) {
			this.install_mode = 'Mac'
		} else if (app.indexOf('Linux') != -1) {
			this.install_mode = 'Linux'
		} else {
			this.install_mode = 'anywhere with Go compiler'
		}
	},
	computed: {
		other_install_modes() {
			let current = this.install_mode
			return ['Windows', 'Mac', 'Linux', 'anywhere with Go compiler'].filter(os => os !== current)
		},
	},
	components: {
		Banner,
		InstallWindows,
		InstallMac,
		InstallLinux,
		InstallGo,
	}
}
</script>