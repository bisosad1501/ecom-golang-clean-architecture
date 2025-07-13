'use client'

import React, { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Progress } from '@/components/ui/progress'
import { Separator } from '@/components/ui/separator'
import { 
  Search, 
  Globe, 
  TrendingUp, 
  AlertTriangle, 
  CheckCircle, 
  Lightbulb,
  ExternalLink,
  Copy,
  RefreshCw,
  BarChart3,
  Users,
  Target
} from 'lucide-react'
import { useCategorySEO } from '@/hooks/use-categories'
import { toast } from 'sonner'

interface CategorySEOPageProps {
  categoryId: string
}

export default function AdminCategorySEOPage({ categoryId }: CategorySEOPageProps) {
  const router = useRouter()
  const [activeTab, setActiveTab] = useState('overview')
  const [isLoading, setIsLoading] = useState(false)

  const {
    categorySEO,
    seoInsights,
    competitorAnalysis,
    slugSuggestions,
    isLoadingSEO,
    updateCategorySEO,
    generateSEO,
    validateSEO,
    optimizeSlug,
    generateSlugSuggestions
  } = useCategorySEO(categoryId)

  const [seoForm, setSeoForm] = useState({
    meta_title: '',
    meta_description: '',
    meta_keywords: '',
    canonical_url: '',
    og_title: '',
    og_description: '',
    og_image: '',
    twitter_title: '',
    twitter_description: '',
    twitter_image: '',
    schema_markup: ''
  })

  useEffect(() => {
    if (categorySEO) {
      setSeoForm({
        meta_title: categorySEO.meta_title || '',
        meta_description: categorySEO.meta_description || '',
        meta_keywords: categorySEO.meta_keywords || '',
        canonical_url: categorySEO.canonical_url || '',
        og_title: categorySEO.og_title || '',
        og_description: categorySEO.og_description || '',
        og_image: categorySEO.og_image || '',
        twitter_title: categorySEO.twitter_title || '',
        twitter_description: categorySEO.twitter_description || '',
        twitter_image: categorySEO.twitter_image || '',
        schema_markup: categorySEO.schema_markup || ''
      })
    }
  }, [categorySEO])

  const handleSaveSEO = async () => {
    setIsLoading(true)
    try {
      await updateCategorySEO(seoForm)
      toast.success('SEO metadata updated successfully')
    } catch (error) {
      toast.error('Failed to update SEO metadata')
    } finally {
      setIsLoading(false)
    }
  }

  const handleGenerateSEO = async () => {
    setIsLoading(true)
    try {
      const generated = await generateSEO()
      setSeoForm(prev => ({
        ...prev,
        ...generated
      }))
      toast.success('SEO metadata generated successfully')
    } catch (error) {
      toast.error('Failed to generate SEO metadata')
    } finally {
      setIsLoading(false)
    }
  }

  const handleOptimizeSlug = async (newSlug: string) => {
    setIsLoading(true)
    try {
      await optimizeSlug({
        new_slug: newSlug,
        preserve_history: true,
        auto_redirect: true
      })
      toast.success('Slug optimized successfully')
    } catch (error) {
      toast.error('Failed to optimize slug')
    } finally {
      setIsLoading(false)
    }
  }

  const getSEOGrade = (score: number) => {
    if (score >= 90) return { grade: 'A+', color: 'bg-green-500' }
    if (score >= 80) return { grade: 'A', color: 'bg-green-400' }
    if (score >= 70) return { grade: 'B', color: 'bg-yellow-500' }
    if (score >= 60) return { grade: 'C', color: 'bg-orange-500' }
    if (score >= 50) return { grade: 'D', color: 'bg-red-400' }
    return { grade: 'F', color: 'bg-red-500' }
  }

  if (isLoadingSEO) {
    return (
      <div className="flex items-center justify-center h-64">
        <RefreshCw className="h-8 w-8 animate-spin" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Category SEO Management</h1>
          <p className="text-muted-foreground">
            Optimize your category for search engines and social media
          </p>
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="outline" onClick={handleGenerateSEO} disabled={isLoading}>
            <Lightbulb className="h-4 w-4 mr-2" />
            Auto Generate
          </Button>
          <Button onClick={handleSaveSEO} disabled={isLoading}>
            Save Changes
          </Button>
        </div>
      </div>

      {/* SEO Score Overview */}
      {seoInsights && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center space-x-2">
              <BarChart3 className="h-5 w-5" />
              <span>SEO Performance</span>
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
              <div className="text-center">
                <div className={`inline-flex items-center justify-center w-16 h-16 rounded-full text-white text-xl font-bold ${getSEOGrade(seoInsights.current_seo.score).color}`}>
                  {getSEOGrade(seoInsights.current_seo.score).grade}
                </div>
                <p className="mt-2 text-sm text-muted-foreground">SEO Grade</p>
                <p className="text-2xl font-bold">{seoInsights.current_seo.score}/100</p>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-red-500">{seoInsights.current_seo.issues.length}</div>
                <p className="text-sm text-muted-foreground">Issues Found</p>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-blue-500">{seoInsights.current_seo.suggestions.length}</div>
                <p className="text-sm text-muted-foreground">Suggestions</p>
              </div>
              <div className="text-center">
                <div className="text-2xl font-bold text-green-500">{seoInsights.recommendations.priority.length}</div>
                <p className="text-sm text-muted-foreground">Priority Actions</p>
              </div>
            </div>
          </CardContent>
        </Card>
      )}

      {/* Main Content Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-6">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="metadata">Metadata</TabsTrigger>
          <TabsTrigger value="social">Social Media</TabsTrigger>
          <TabsTrigger value="technical">Technical</TabsTrigger>
          <TabsTrigger value="insights">Insights</TabsTrigger>
          <TabsTrigger value="competitors">Competitors</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-6">
          {/* Quick Actions */}
          <Card>
            <CardHeader>
              <CardTitle>Quick Actions</CardTitle>
              <CardDescription>
                Common SEO optimization tasks
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <Button variant="outline" className="h-auto p-4 flex flex-col items-start">
                  <Search className="h-5 w-5 mb-2" />
                  <span className="font-medium">Optimize Meta Tags</span>
                  <span className="text-sm text-muted-foreground">Improve title and description</span>
                </Button>
                <Button variant="outline" className="h-auto p-4 flex flex-col items-start">
                  <Globe className="h-5 w-5 mb-2" />
                  <span className="font-medium">Generate Schema</span>
                  <span className="text-sm text-muted-foreground">Add structured data</span>
                </Button>
                <Button variant="outline" className="h-auto p-4 flex flex-col items-start">
                  <TrendingUp className="h-5 w-5 mb-2" />
                  <span className="font-medium">Analyze Performance</span>
                  <span className="text-sm text-muted-foreground">View SEO metrics</span>
                </Button>
              </div>
            </CardContent>
          </Card>

          {/* Issues and Suggestions */}
          {seoInsights && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center space-x-2">
                    <AlertTriangle className="h-5 w-5 text-red-500" />
                    <span>Issues ({seoInsights.current_seo.issues.length})</span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {seoInsights.current_seo.issues.map((issue, index) => (
                      <Alert key={index} variant={issue.severity === 'error' ? 'destructive' : 'default'}>
                        <AlertDescription>
                          <strong>{issue.field}:</strong> {issue.description}
                        </AlertDescription>
                      </Alert>
                    ))}
                    {seoInsights.current_seo.issues.length === 0 && (
                      <div className="text-center py-4 text-muted-foreground">
                        <CheckCircle className="h-8 w-8 mx-auto mb-2 text-green-500" />
                        No issues found
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center space-x-2">
                    <Lightbulb className="h-5 w-5 text-blue-500" />
                    <span>Suggestions ({seoInsights.current_seo.suggestions.length})</span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    {seoInsights.current_seo.suggestions.map((suggestion, index) => (
                      <Alert key={index}>
                        <AlertDescription>
                          <strong>{suggestion.field}:</strong> {suggestion.suggestion}
                          <Badge variant="outline" className="ml-2">
                            {suggestion.impact} impact
                          </Badge>
                        </AlertDescription>
                      </Alert>
                    ))}
                    {seoInsights.current_seo.suggestions.length === 0 && (
                      <div className="text-center py-4 text-muted-foreground">
                        <CheckCircle className="h-8 w-8 mx-auto mb-2 text-green-500" />
                        No suggestions available
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            </div>
          )}
        </TabsContent>

        <TabsContent value="metadata" className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Basic SEO Metadata</CardTitle>
              <CardDescription>
                Essential meta tags for search engine optimization
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="meta_title">Meta Title</Label>
                <Input
                  id="meta_title"
                  value={seoForm.meta_title}
                  onChange={(e) => setSeoForm(prev => ({ ...prev, meta_title: e.target.value }))}
                  placeholder="Enter meta title (50-60 characters)"
                  maxLength={60}
                />
                <p className="text-sm text-muted-foreground">
                  {seoForm.meta_title.length}/60 characters
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="meta_description">Meta Description</Label>
                <Textarea
                  id="meta_description"
                  value={seoForm.meta_description}
                  onChange={(e) => setSeoForm(prev => ({ ...prev, meta_description: e.target.value }))}
                  placeholder="Enter meta description (120-160 characters)"
                  maxLength={160}
                  rows={3}
                />
                <p className="text-sm text-muted-foreground">
                  {seoForm.meta_description.length}/160 characters
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="meta_keywords">Meta Keywords</Label>
                <Input
                  id="meta_keywords"
                  value={seoForm.meta_keywords}
                  onChange={(e) => setSeoForm(prev => ({ ...prev, meta_keywords: e.target.value }))}
                  placeholder="Enter keywords separated by commas"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="canonical_url">Canonical URL</Label>
                <Input
                  id="canonical_url"
                  value={seoForm.canonical_url}
                  onChange={(e) => setSeoForm(prev => ({ ...prev, canonical_url: e.target.value }))}
                  placeholder="Enter canonical URL"
                />
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        {/* Additional tabs would be implemented here */}
      </Tabs>
    </div>
  )
}
