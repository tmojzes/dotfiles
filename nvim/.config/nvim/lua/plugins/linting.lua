return {
	{
		"mfussenegger/nvim-lint",
		opts = {
			-- Event to trigger linters
			linters_by_ft = {
				["*"] = { "codespell" },
			},
		},
	},
}
