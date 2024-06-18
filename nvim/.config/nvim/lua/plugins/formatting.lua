return {
	{
		"stevearc/conform.nvim",
		opts = {
			formatters_by_ft = {
				python = { "isort", "black" },
				lua = { "stylua" },
				go = { "gofmt", "goimports", "golines" },
				sh = { "shfmt" },
			},
		},
	},
}
