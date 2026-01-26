return {
	"leoluz/nvim-dap-go",
	opts = {
		-- Dieser Befehl sagt dem Go-Plugin, wo es dlv finden soll,
		-- und ignoriert das PATH-Problem.
		dap_bin = vim.fn.stdpath("data") .. "/mason/bin/dlv",
	},
}
