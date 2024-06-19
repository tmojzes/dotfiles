return {
	{
		"stevearc/conform.nvim",
		opts = {
			formatters = {
				golines = {
					prepend_args = { "-m", "110" },
				},
			},
			formatters_by_ft = {
				python = { "isort", "black" },
				lua = { "stylua" },
				go = { "gofmt", "goimports", "golines" },
				sh = { "shfmt" },
			},
		},
	},
}
